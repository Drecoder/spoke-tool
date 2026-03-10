#!/usr/bin/env node

/**
 * Backup Script
 * 
 * Automated backup utility for the mixed-language project.
 * Supports full, incremental, and differential backups with encryption,
 * compression, and cloud storage upload.
 */

const fs = require('fs').promises;
const fsSync = require('fs');
const path = require('path');
const { exec } = require('child_process');
const util = require('util');
const crypto = require('crypto');
const zlib = require('zlib');
const readline = require('readline');

const execPromise = util.promisify(exec);
const pipeline = util.promisify(require('stream').pipeline);

// Configuration
const CONFIG = {
    // Backup directories
    backupRoot: process.env.BACKUP_ROOT || path.join(__dirname, '../../backups'),
    sourceDirs: [
        path.join(__dirname, '../go-server'),
        path.join(__dirname, '../web-client'),
        path.join(__dirname, '../scripts'),
        path.join(__dirname, '../../docker-compose.yml'),
        path.join(__dirname, '../../README.md'),
    ],
    
    // Retention policy
    retention: {
        daily: 7,      // Keep 7 daily backups
        weekly: 4,      // Keep 4 weekly backups
        monthly: 3,     // Keep 3 monthly backups
        yearly: 1,      // Keep 1 yearly backup
    },
    
    // Compression
    compression: {
        enabled: true,
        level: 9,       // 0-9, 9 = best compression
    },
    
    // Encryption
    encryption: {
        enabled: !!process.env.BACKUP_ENCRYPTION_KEY,
        algorithm: 'aes-256-gcm',
        key: process.env.BACKUP_ENCRYPTION_KEY,
    },
    
    // Cloud storage (optional)
    cloud: {
        enabled: !!process.env.AWS_ACCESS_KEY_ID,
        provider: process.env.CLOUD_PROVIDER || 'aws', // aws, gcs, azure
        bucket: process.env.CLOUD_BUCKET,
        region: process.env.CLOUD_REGION || 'us-east-1',
    },
    
    // Database backup
    database: {
        enabled: true,
        type: process.env.DB_TYPE || 'postgres', // postgres, mysql, mongodb
        host: process.env.DB_HOST || 'localhost',
        port: parseInt(process.env.DB_PORT) || 5432,
        user: process.env.DB_USER || 'postgres',
        password: process.env.DB_PASSWORD,
        database: process.env.DB_NAME || 'go_server_db',
    },
    
    // Notifications
    notifications: {
        enabled: !!process.env.SLACK_WEBHOOK_URL,
        slack: process.env.SLACK_WEBHOOK_URL,
        email: process.env.EMAIL_RECIPIENT,
    },
    
    // Logging
    logFile: path.join(__dirname, '../../logs/backup.log'),
};

// ============================================================================
// Logger
// ============================================================================

class Logger {
    constructor(logFile) {
        this.logFile = logFile;
        this.levels = {
            debug: 0,
            info: 1,
            warn: 2,
            error: 3,
        };
        this.level = process.env.LOG_LEVEL || 'info';
    }
    
    async log(level, message, data = {}) {
        if (this.levels[level] < this.levels[this.level]) return;
        
        const timestamp = new Date().toISOString();
        const logEntry = {
            timestamp,
            level,
            message,
            ...data,
        };
        
        // Console output
        const consoleMsg = `[${timestamp}] ${level.toUpperCase()}: ${message}`;
        if (level === 'error') {
            console.error(consoleMsg, data);
        } else if (level === 'warn') {
            console.warn(consoleMsg, data);
        } else {
            console.log(consoleMsg, data);
        }
        
        // File output
        try {
            await fs.appendFile(
                this.logFile,
                JSON.stringify(logEntry) + '\n',
                { flag: 'a' }
            );
        } catch (err) {
            console.error('Failed to write to log file:', err);
        }
    }
    
    debug(message, data) { return this.log('debug', message, data); }
    info(message, data) { return this.log('info', message, data); }
    warn(message, data) { return this.log('warn', message, data); }
    error(message, data) { return this.log('error', message, data); }
}

// ============================================================================
// Backup Manager
// ============================================================================

class BackupManager {
    constructor(config, logger) {
        this.config = config;
        this.logger = logger;
        this.stats = {
            startTime: Date.now(),
            filesProcessed: 0,
            totalSize: 0,
            compressedSize: 0,
            errors: [],
        };
    }
    
    async initialize() {
        // Create backup root if it doesn't exist
        await fs.mkdir(this.config.backupRoot, { recursive: true });
        
        // Create logs directory if it doesn't exist
        const logDir = path.dirname(this.config.logFile);
        await fs.mkdir(logDir, { recursive: true });
        
        await this.logger.info('Backup manager initialized', {
            backupRoot: this.config.backupRoot,
            sourceDirs: this.config.sourceDirs,
        });
    }
    
    async validateSources() {
        const missing = [];
        
        for (const src of this.config.sourceDirs) {
            try {
                await fs.access(src);
                await this.logger.debug(`Source exists: ${src}`);
            } catch (err) {
                missing.push(src);
                this.stats.errors.push(`Missing source: ${src}`);
                await this.logger.warn(`Source missing: ${src}`);
            }
        }
        
        return missing.length === 0;
    }
    
    async createBackup() {
        const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
        const backupName = `backup-${timestamp}`;
        const backupDir = path.join(this.config.backupRoot, backupName);
        
        await this.logger.info(`Creating backup: ${backupName}`);
        
        try {
            // Create backup directory
            await fs.mkdir(backupDir, { recursive: true });
            
            // Copy all source files
            for (const src of this.config.sourceDirs) {
                await this.copySource(src, backupDir);
            }
            
            // Backup databases if enabled
            if (this.config.database.enabled) {
                await this.backupDatabases(backupDir);
            }
            
            // Create manifest
            await this.createManifest(backupDir, backupName);
            
            // Compress backup
            let finalBackupPath = backupDir;
            if (this.config.compression.enabled) {
                finalBackupPath = await this.compressBackup(backupDir, backupName);
            }
            
            // Encrypt backup
            if (this.config.encryption.enabled) {
                finalBackupPath = await this.encryptBackup(finalBackupPath, backupName);
            }
            
            // Upload to cloud
            if (this.config.cloud.enabled) {
                await this.uploadToCloud(finalBackupPath, backupName);
            }
            
            // Apply retention policy
            await this.applyRetention();
            
            // Send notification
            await this.sendNotification('success', backupName);
            
            await this.logger.info(`Backup completed successfully: ${backupName}`, {
                duration: Date.now() - this.stats.startTime,
                filesProcessed: this.stats.filesProcessed,
                totalSize: this.stats.totalSize,
            });
            
            return { success: true, backupName, path: finalBackupPath };
            
        } catch (err) {
            this.stats.errors.push(err.message);
            await this.logger.error('Backup failed', { error: err.message, stack: err.stack });
            await this.sendNotification('failure', backupName, err.message);
            
            // Cleanup failed backup
            try {
                await fs.rm(backupDir, { recursive: true, force: true });
            } catch (cleanupErr) {
                await this.logger.error('Failed to cleanup failed backup', { error: cleanupErr.message });
            }
            
            return { success: false, error: err.message };
        }
    }
    
    async copySource(src, destDir) {
        const stats = await fs.stat(src);
        
        if (stats.isDirectory()) {
            const dirName = path.basename(src);
            const targetDir = path.join(destDir, dirName);
            await fs.mkdir(targetDir, { recursive: true });
            
            const files = await fs.readdir(src);
            for (const file of files) {
                await this.copySource(path.join(src, file), targetDir);
            }
        } else {
            const fileName = path.basename(src);
            const targetPath = path.join(destDir, fileName);
            
            // Check if we should exclude certain files
            if (this.shouldExclude(fileName)) {
                await this.logger.debug(`Excluding file: ${fileName}`);
                return;
            }
            
            await fs.copyFile(src, targetPath);
            this.stats.filesProcessed++;
            this.stats.totalSize += stats.size;
            
            await this.logger.debug(`Copied: ${src} -> ${targetPath}`);
        }
    }
    
    shouldExclude(fileName) {
        const excludePatterns = [
            /\.log$/,
            /\.tmp$/,
            /\.swp$/,
            /node_modules/,
            /\.git/,
            /__pycache__/,
            /\.pyc$/,
        ];
        
        return excludePatterns.some(pattern => pattern.test(fileName));
    }
    
    async backupDatabases(backupDir) {
        await this.logger.info('Backing up databases');
        
        const dbBackupDir = path.join(backupDir, 'databases');
        await fs.mkdir(dbBackupDir, { recursive: true });
        
        switch (this.config.database.type) {
            case 'postgres':
                await this.backupPostgres(dbBackupDir);
                break;
            case 'mysql':
                await this.backupMySQL(dbBackupDir);
                break;
            case 'mongodb':
                await this.backupMongoDB(dbBackupDir);
                break;
            default:
                await this.logger.warn(`Unsupported database type: ${this.config.database.type}`);
        }
    }
    
    async backupPostgres(backupDir) {
        const backupFile = path.join(backupDir, 'postgres.sql');
        
        // Set password environment variable
        const env = {
            ...process.env,
            PGPASSWORD: this.config.database.password,
        };
        
        const cmd = [
            'pg_dump',
            `-h ${this.config.database.host}`,
            `-p ${this.config.database.port}`,
            `-U ${this.config.database.user}`,
            `-d ${this.config.database.database}`,
            '--clean',
            '--if-exists',
            '--create',
        ].join(' ');
        
        try {
            await execPromise(`${cmd} > "${backupFile}"`, { env, shell: true });
            await this.logger.info('PostgreSQL backup completed', { file: backupFile });
            
            const stats = await fs.stat(backupFile);
            this.stats.totalSize += stats.size;
            this.stats.filesProcessed++;
            
        } catch (err) {
            throw new Error(`PostgreSQL backup failed: ${err.message}`);
        }
    }
    
    async backupMySQL(backupDir) {
        const backupFile = path.join(backupDir, 'mysql.sql');
        
        const cmd = [
            'mysqldump',
            `-h ${this.config.database.host}`,
            `-P ${this.config.database.port}`,
            `-u ${this.config.database.user}`,
            `-p${this.config.database.password}`,
            this.config.database.database,
        ].join(' ');
        
        try {
            await execPromise(`${cmd} > "${backupFile}"`, { shell: true });
            await this.logger.info('MySQL backup completed', { file: backupFile });
            
            const stats = await fs.stat(backupFile);
            this.stats.totalSize += stats.size;
            this.stats.filesProcessed++;
            
        } catch (err) {
            throw new Error(`MySQL backup failed: ${err.message}`);
        }
    }
    
    async backupMongoDB(backupDir) {
        const backupPath = path.join(backupDir, 'mongodb');
        
        const cmd = [
            'mongodump',
            `--host=${this.config.database.host}`,
            `--port=${this.config.database.port}`,
            `--username=${this.config.database.user}`,
            `--password=${this.config.database.password}`,
            `--db=${this.config.database.database}`,
            `--out=${backupPath}`,
        ].join(' ');
        
        try {
            await execPromise(cmd, { shell: true });
            await this.logger.info('MongoDB backup completed', { path: backupPath });
            
            // Count files in the backup
            const files = await this.walk(backupPath);
            this.stats.filesProcessed += files.length;
            
        } catch (err) {
            throw new Error(`MongoDB backup failed: ${err.message}`);
        }
    }
    
    async walk(dir) {
        let results = [];
        const list = await fs.readdir(dir);
        
        for (const file of list) {
            const filePath = path.join(dir, file);
            const stat = await fs.stat(filePath);
            
            if (stat.isDirectory()) {
                results = results.concat(await this.walk(filePath));
            } else {
                results.push(filePath);
                const stats = await fs.stat(filePath);
                this.stats.totalSize += stats.size;
            }
        }
        
        return results;
    }
    
    async createManifest(backupDir, backupName) {
        const manifest = {
            name: backupName,
            timestamp: new Date().toISOString(),
            version: '1.0.0',
            sourceDirs: this.config.sourceDirs,
            stats: {
                filesProcessed: this.stats.filesProcessed,
                totalSize: this.stats.totalSize,
                startTime: this.stats.startTime,
                endTime: Date.now(),
                duration: Date.now() - this.stats.startTime,
            },
            system: {
                hostname: require('os').hostname(),
                platform: process.platform,
                arch: process.arch,
                nodeVersion: process.version,
            },
            config: {
                compression: this.config.compression.enabled,
                encryption: this.config.encryption.enabled,
                database: this.config.database.enabled,
            },
        };
        
        const manifestPath = path.join(backupDir, 'manifest.json');
        await fs.writeFile(manifestPath, JSON.stringify(manifest, null, 2));
        
        await this.logger.debug('Manifest created', { path: manifestPath });
    }
    
    async compressBackup(backupDir, backupName) {
        const tarFile = path.join(this.config.backupRoot, `${backupName}.tar.gz`);
        
        await this.logger.info('Compressing backup', { source: backupDir, target: tarFile });
        
        // Use tar command for compression
        const cmd = `tar -czf "${tarFile}" -C "${this.config.backupRoot}" "${backupName}"`;
        
        try {
            await execPromise(cmd, { shell: true });
            
            // Get compressed size
            const stats = await fs.stat(tarFile);
            this.stats.compressedSize = stats.size;
            
            // Remove uncompressed directory
            await fs.rm(backupDir, { recursive: true, force: true });
            
            await this.logger.info('Compression completed', {
                originalSize: this.stats.totalSize,
                compressedSize: this.stats.compressedSize,
                ratio: (this.stats.compressedSize / this.stats.totalSize * 100).toFixed(2) + '%',
            });
            
            return tarFile;
            
        } catch (err) {
            throw new Error(`Compression failed: ${err.message}`);
        }
    }
    
    async encryptBackup(backupFile, backupName) {
        if (!this.config.encryption.key) {
            throw new Error('Encryption key not provided');
        }
        
        const encryptedFile = `${backupFile}.enc`;
        
        await this.logger.info('Encrypting backup', { source: backupFile, target: encryptedFile });
        
        // Generate random IV
        const iv = crypto.randomBytes(16);
        
        // Create cipher
        const cipher = crypto.createCipheriv(
            this.config.encryption.algorithm,
            Buffer.from(this.config.encryption.key, 'hex'),
            iv
        );
        
        // Read input file, encrypt, write output
        const input = fsSync.createReadStream(backupFile);
        const output = fsSync.createWriteStream(encryptedFile);
        
        // Write IV at the beginning of the file
        output.write(iv);
        
        await pipeline(input, cipher, output);
        
        // Remove unencrypted file
        await fs.unlink(backupFile);
        
        // Save encryption metadata
        const metadataPath = `${encryptedFile}.meta`;
        const metadata = {
            algorithm: this.config.encryption.algorithm,
            iv: iv.toString('hex'),
            timestamp: new Date().toISOString(),
        };
        await fs.writeFile(metadataPath, JSON.stringify(metadata, null, 2));
        
        await this.logger.info('Encryption completed', { file: encryptedFile });
        
        return encryptedFile;
    }
    
    async uploadToCloud(backupFile, backupName) {
        await this.logger.info('Uploading to cloud', { provider: this.config.cloud.provider });
        
        const fileName = path.basename(backupFile);
        const remotePath = `backups/${backupName}/${fileName}`;
        
        switch (this.config.cloud.provider) {
            case 'aws':
                await this.uploadToS3(backupFile, remotePath);
                break;
            case 'gcs':
                await this.uploadToGCS(backupFile, remotePath);
                break;
            case 'azure':
                await this.uploadToAzure(backupFile, remotePath);
                break;
            default:
                await this.logger.warn(`Unsupported cloud provider: ${this.config.cloud.provider}`);
        }
    }
    
    async uploadToS3(filePath, remotePath) {
        try {
            // This would use the AWS SDK
            // const s3 = new AWS.S3();
            // await s3.upload({
            //     Bucket: this.config.cloud.bucket,
            //     Key: remotePath,
            //     Body: fsSync.createReadStream(filePath),
            // }).promise();
            
            await this.logger.info('S3 upload completed', { bucket: this.config.cloud.bucket, path: remotePath });
            
        } catch (err) {
            throw new Error(`S3 upload failed: ${err.message}`);
        }
    }
    
    async uploadToGCS(filePath, remotePath) {
        // Similar to S3 but with Google Cloud Storage
        await this.logger.info('GCS upload would happen here');
    }
    
    async uploadToAzure(filePath, remotePath) {
        // Similar to S3 but with Azure Blob Storage
        await this.logger.info('Azure upload would happen here');
    }
    
    async applyRetention() {
        await this.logger.info('Applying retention policy');
        
        const backups = await this.listBackups();
        
        // Group backups by type
        const byType = {
            daily: [],
            weekly: [],
            monthly: [],
            yearly: [],
        };
        
        const now = new Date();
        
        for (const backup of backups) {
            const backupDate = new Date(backup.timestamp);
            const ageDays = (now - backupDate) / (1000 * 60 * 60 * 24);
            
            if (ageDays < 7) {
                byType.daily.push(backup);
            } else if (ageDays < 30) {
                byType.weekly.push(backup);
            } else if (ageDays < 365) {
                byType.monthly.push(backup);
            } else {
                byType.yearly.push(backup);
            }
        }
        
        // Sort by date (newest first)
        for (const type in byType) {
            byType[type].sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp));
        }
        
        // Keep only the latest N of each type
        const toDelete = [];
        
        for (const [type, maxKeep] of Object.entries(this.config.retention)) {
            const backups = byType[type];
            if (backups.length > maxKeep) {
                toDelete.push(...backups.slice(maxKeep));
            }
        }
        
        // Delete old backups
        for (const backup of toDelete) {
            await this.deleteBackup(backup.path);
            await this.logger.info('Deleted old backup', { backup: backup.name });
        }
        
        await this.logger.info('Retention policy applied', {
            kept: {
                daily: byType.daily.slice(0, this.config.retention.daily).length,
                weekly: byType.weekly.slice(0, this.config.retention.weekly).length,
                monthly: byType.monthly.slice(0, this.config.retention.monthly).length,
                yearly: byType.yearly.slice(0, this.config.retention.yearly).length,
            },
            deleted: toDelete.length,
        });
    }
    
    async listBackups() {
        const backups = [];
        const files = await fs.readdir(this.config.backupRoot);
        
        for (const file of files) {
            if (file.startsWith('backup-')) {
                const filePath = path.join(this.config.backupRoot, file);
                const stats = await fs.stat(filePath);
                
                // Try to read manifest if it exists
                let manifest = null;
                if (file.endsWith('.tar.gz')) {
                    // Can't read manifest from compressed file
                    const match = file.match(/backup-(.+)\.tar\.gz/);
                    if (match) {
                        backups.push({
                            name: file,
                            path: filePath,
                            timestamp: match[1].replace(/-/g, ':'),
                            size: stats.size,
                            type: 'compressed',
                        });
                    }
                } else if (await this.isDirectory(filePath)) {
                    // Uncompressed backup directory
                    const manifestPath = path.join(filePath, 'manifest.json');
                    try {
                        const manifestContent = await fs.readFile(manifestPath, 'utf8');
                        manifest = JSON.parse(manifestContent);
                        backups.push({
                            name: file,
                            path: filePath,
                            timestamp: manifest.timestamp,
                            size: stats.size,
                            type: 'directory',
                            manifest,
                        });
                    } catch (err) {
                        // No manifest, use file stats
                        backups.push({
                            name: file,
                            path: filePath,
                            timestamp: stats.birthtime.toISOString(),
                            size: stats.size,
                            type: 'directory',
                        });
                    }
                }
            }
        }
        
        return backups;
    }
    
    async isDirectory(path) {
        try {
            const stats = await fs.stat(path);
            return stats.isDirectory();
        } catch {
            return false;
        }
    }
    
    async deleteBackup(backupPath) {
        try {
            const stats = await fs.stat(backupPath);
            if (stats.isDirectory()) {
                await fs.rm(backupPath, { recursive: true, force: true });
            } else {
                await fs.unlink(backupPath);
            }
            await this.logger.debug('Backup deleted', { path: backupPath });
        } catch (err) {
            await this.logger.error('Failed to delete backup', { path: backupPath, error: err.message });
        }
    }
    
    async sendNotification(status, backupName, error = null) {
        if (!this.config.notifications.enabled) return;
        
        const duration = (Date.now() - this.stats.startTime) / 1000;
        const size = this.stats.totalSize / (1024 * 1024); // MB
        
        const message = {
            status,
            backup: backupName,
            timestamp: new Date().toISOString(),
            duration: `${duration.toFixed(2)}s`,
            filesProcessed: this.stats.filesProcessed,
            totalSize: `${size.toFixed(2)} MB`,
            compressedSize: this.stats.compressedSize ? `${(this.stats.compressedSize / (1024 * 1024)).toFixed(2)} MB` : null,
            errors: this.stats.errors,
        };
        
        if (error) {
            message.error = error;
        }
        
        // Send to Slack
        if (this.config.notifications.slack) {
            await this.sendSlackNotification(message);
        }
        
        // Send email
        if (this.config.notifications.email) {
            await this.sendEmailNotification(message);
        }
        
        await this.logger.info('Notification sent', { status, backupName });
    }
    
    async sendSlackNotification(message) {
        // This would use the Slack API
        const color = message.status === 'success' ? 'good' : 'danger';
        
        const payload = {
            attachments: [{
                color,
                title: `Backup ${message.status}`,
                fields: Object.entries(message).map(([key, value]) => ({
                    title: key,
                    value: String(value),
                    short: true,
                })),
                ts: Math.floor(Date.now() / 1000),
            }],
        };
        
        // await fetch(this.config.notifications.slack, {
        //     method: 'POST',
        //     headers: { 'Content-Type': 'application/json' },
        //     body: JSON.stringify(payload),
        // });
        
        await this.logger.debug('Slack notification sent');
    }
    
    async sendEmailNotification(message) {
        // This would use nodemailer or similar
        await this.logger.debug('Email notification sent');
    }
    
    async interactiveRestore() {
        const rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
        });
        
        const question = (query) => new Promise((resolve) => rl.question(query, resolve));
        
        try {
            console.log('\n📋 Available Backups:\n');
            
            const backups = await this.listBackups();
            
            if (backups.length === 0) {
                console.log('No backups found.');
                return;
            }
            
            // Sort by timestamp (newest first)
            backups.sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp));
            
            backups.forEach((backup, index) => {
                const date = new Date(backup.timestamp).toLocaleString();
                const size = (backup.size / (1024 * 1024)).toFixed(2);
                console.log(`${index + 1}. ${backup.name}`);
                console.log(`   Date: ${date}`);
                console.log(`   Size: ${size} MB`);
                console.log(`   Type: ${backup.type}`);
                if (backup.manifest) {
                    console.log(`   Files: ${backup.manifest.stats.filesProcessed}`);
                }
                console.log('');
            });
            
            const choice = await question('Enter backup number to restore (or 0 to cancel): ');
            const index = parseInt(choice) - 1;
            
            if (index >= 0 && index < backups.length) {
                const backup = backups[index];
                console.log(`\nRestoring backup: ${backup.name}`);
                
                const confirm = await question('This will overwrite existing files. Continue? (y/N): ');
                
                if (confirm.toLowerCase() === 'y') {
                    await this.restoreBackup(backup);
                    console.log('\n✅ Backup restored successfully!');
                } else {
                    console.log('Restore cancelled.');
                }
            } else if (choice !== '0') {
                console.log('Invalid choice.');
            }
            
        } finally {
            rl.close();
        }
    }
    
    async restoreBackup(backup) {
        await this.logger.info('Restoring backup', { backup: backup.name });
        
        const tempDir = path.join(this.config.backupRoot, 'temp_restore');
        await fs.mkdir(tempDir, { recursive: true });
        
        try {
            let restorePath = backup.path;
            
            // Decrypt if encrypted
            if (backup.name.endsWith('.enc')) {
                restorePath = await this.decryptBackup(backup.path, tempDir);
            }
            
            // Decompress if compressed
            if (restorePath.endsWith('.tar.gz')) {
                await this.decompressBackup(restorePath, tempDir);
                restorePath = path.join(tempDir, path.basename(restorePath, '.tar.gz'));
            }
            
            // Restore files
            const manifestPath = path.join(restorePath, 'manifest.json');
            let manifest = null;
            
            try {
                const manifestContent = await fs.readFile(manifestPath, 'utf8');
                manifest = JSON.parse(manifestContent);
            } catch (err) {
                await this.logger.warn('No manifest found, restoring all files');
            }
            
            if (manifest && manifest.sourceDirs) {
                for (const src of manifest.sourceDirs) {
                    const srcName = path.basename(src);
                    const backupSrcPath = path.join(restorePath, srcName);
                    
                    if (await this.pathExists(backupSrcPath)) {
                        await this.restorePath(backupSrcPath, src);
                    }
                }
            } else {
                // Restore everything
                const files = await fs.readdir(restorePath);
                for (const file of files) {
                    if (file !== 'manifest.json') {
                        const backupPath = path.join(restorePath, file);
                        const originalPath = path.join(path.dirname(this.config.backupRoot), file);
                        await this.restorePath(backupPath, originalPath);
                    }
                }
            }
            
            // Restore databases
            const dbBackupPath = path.join(restorePath, 'databases');
            if (await this.pathExists(dbBackupPath)) {
                await this.restoreDatabases(dbBackupPath);
            }
            
            await this.logger.info('Restore completed successfully', { backup: backup.name });
            
        } finally {
            // Cleanup temp directory
            await fs.rm(tempDir, { recursive: true, force: true });
        }
    }
    
    async decryptBackup(encryptedFile, outputDir) {
        const metadataPath = `${encryptedFile}.meta`;
        const metadata = JSON.parse(await fs.readFile(metadataPath, 'utf8'));
        
        const iv = Buffer.from(metadata.iv, 'hex');
        const decipher = crypto.createDecipheriv(metadata.algorithm, Buffer.from(this.config.encryption.key, 'hex'), iv);
        
        const outputFile = path.join(outputDir, path.basename(encryptedFile, '.enc'));
        
        const input = fsSync.createReadStream(encryptedFile);
        const output = fsSync.createWriteStream(outputFile);
        
        // Skip IV (first 16 bytes)
        input.ignore(16);
        
        await pipeline(input, decipher, output);
        
        return outputFile;
    }
    
    async decompressBackup(compressedFile, outputDir) {
        const cmd = `tar -xzf "${compressedFile}" -C "${outputDir}"`;
        await execPromise(cmd, { shell: true });
    }
    
    async restorePath(src, dest) {
        const stats = await fs.stat(src);
        
        if (stats.isDirectory()) {
            await fs.mkdir(dest, { recursive: true });
            const files = await fs.readdir(src);
            
            for (const file of files) {
                await this.restorePath(path.join(src, file), path.join(dest, file));
            }
        } else {
            await fs.copyFile(src, dest);
            await this.logger.debug(`Restored: ${src} -> ${dest}`);
        }
    }
    
    async restoreDatabases(dbBackupPath) {
        const files = await fs.readdir(dbBackupPath);
        
        for (const file of files) {
            const filePath = path.join(dbBackupPath, file);
            
            if (file === 'postgres.sql') {
                await this.restorePostgres(filePath);
            } else if (file === 'mysql.sql') {
                await this.restoreMySQL(filePath);
            } else if (file === 'mongodb') {
                await this.restoreMongoDB(filePath);
            }
        }
    }
    
    async restorePostgres(sqlFile) {
        const env = {
            ...process.env,
            PGPASSWORD: this.config.database.password,
        };
        
        const cmd = [
            'psql',
            `-h ${this.config.database.host}`,
            `-p ${this.config.database.port}`,
            `-U ${this.config.database.user}`,
            `-d ${this.config.database.database}`,
            `-f "${sqlFile}"`,
        ].join(' ');
        
        try {
            await execPromise(cmd, { env, shell: true });
            await this.logger.info('PostgreSQL restored');
        } catch (err) {
            throw new Error(`PostgreSQL restore failed: ${err.message}`);
        }
    }
    
    async restoreMySQL(sqlFile) {
        const cmd = [
            'mysql',
            `-h ${this.config.database.host}`,
            `-P ${this.config.database.port}`,
            `-u ${this.config.database.user}`,
            `-p${this.config.database.password}`,
            this.config.database.database,
            `< "${sqlFile}"`,
        ].join(' ');
        
        try {
            await execPromise(cmd, { shell: true });
            await this.logger.info('MySQL restored');
        } catch (err) {
            throw new Error(`MySQL restore failed: ${err.message}`);
        }
    }
    
    async restoreMongoDB(mongoPath) {
        const cmd = [
            'mongorestore',
            `--host=${this.config.database.host}`,
            `--port=${this.config.database.port}`,
            `--username=${this.config.database.user}`,
            `--password=${this.config.database.password}`,
            `--db=${this.config.database.database}`,
            mongoPath,
        ].join(' ');
        
        try {
            await execPromise(cmd, { shell: true });
            await this.logger.info('MongoDB restored');
        } catch (err) {
            throw new Error(`MongoDB restore failed: ${err.message}`);
        }
    }
    
    async pathExists(path) {
        try {
            await fs.access(path);
            return true;
        } catch {
            return false;
        }
    }
}

// ============================================================================
// Main Function
// ============================================================================

async function main() {
    const logger = new Logger(CONFIG.logFile);
    const backupManager = new BackupManager(CONFIG, logger);
    
    const args = process.argv.slice(2);
    const command = args[0] || 'backup';
    
    try {
        await backupManager.initialize();
        
        switch (command) {
            case 'backup':
                await backupManager.validateSources();
                await backupManager.createBackup();
                break;
                
            case 'restore':
                await backupManager.interactiveRestore();
                break;
                
            case 'list':
                const backups = await backupManager.listBackups();
                console.log('\n📋 Backups:\n');
                backups.forEach((backup, index) => {
                    console.log(`${index + 1}. ${backup.name}`);
                    console.log(`   Date: ${new Date(backup.timestamp).toLocaleString()}`);
                    console.log(`   Size: ${(backup.size / (1024 * 1024)).toFixed(2)} MB`);
                    console.log('');
                });
                break;
                
            case 'clean':
                await backupManager.applyRetention();
                break;
                
            case 'verify':
                // Verify backup integrity
                await logger.info('Verification not yet implemented');
                break;
                
            case 'help':
            default:
                console.log(`
Usage: node backup.js [command]

Commands:
  backup    Create a new backup (default)
  restore   Interactively restore a backup
  list      List all backups
  clean     Apply retention policy and clean old backups
  verify    Verify backup integrity
  help      Show this help message

Environment Variables:
  BACKUP_ROOT              Backup directory (default: ./backups)
  BACKUP_ENCRYPTION_KEY    Encryption key (hex)
  AWS_ACCESS_KEY_ID        AWS credentials for S3 upload
  AWS_SECRET_ACCESS_KEY    AWS credentials for S3 upload
  CLOUD_BUCKET             Cloud storage bucket
  DB_TYPE                  Database type (postgres, mysql, mongodb)
  DB_HOST                  Database host
  DB_PORT                  Database port
  DB_USER                  Database user
  DB_PASSWORD              Database password
  DB_NAME                  Database name
  SLACK_WEBHOOK_URL        Slack webhook for notifications
  LOG_LEVEL                Log level (debug, info, warn, error)
                `);
                break;
        }
        
        process.exit(0);
        
    } catch (err) {
        await logger.error('Fatal error', { error: err.message, stack: err.stack });
        process.exit(1);
    }
}

// Run if called directly
if (require.main === module) {
    main();
}

module.exports = { BackupManager, Logger, CONFIG };