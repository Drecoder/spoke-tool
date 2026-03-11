#!/usr/bin/env python3
"""
Deployment Script

Automated deployment utility for the mixed-language project.
Supports multiple environments, health checks, rollback, and blue-green deployments.
"""

import os
import sys
import json
import time
import yaml
import argparse
import logging
import subprocess
import requests
import paramiko
from pathlib import Path
from datetime import datetime
from typing import Dict, List, Optional, Tuple, Any
import concurrent.futures
from dataclasses import dataclass, asdict
import hashlib
import hmac

# ============================================================================
# Configuration
# ============================================================================

@dataclass
class ServiceConfig:
    """Service configuration"""
    name: str
    type: str  # 'go', 'nodejs', 'python'
    path: str
    build_command: Optional[str] = None
    start_command: Optional[str] = None
    health_endpoint: Optional[str] = None
    health_timeout: int = 30
    env_vars: Dict[str, str] = None
    port: Optional[int] = None
    replicas: int = 1

@dataclass
class EnvironmentConfig:
    """Environment configuration"""
    name: str
    hosts: List[str]
    ssh_user: str
    ssh_key: str
    deploy_path: str
    backup_path: str
    services: List[ServiceConfig]
    env_vars: Dict[str, str] = None

class DeploymentConfig:
    """Main deployment configuration"""
    
    def __init__(self, config_path: str):
        with open(config_path, 'r') as f:
            if config_path.endswith('.yaml') or config_path.endswith('.yml'):
                self.config = yaml.safe_load(f)
            elif config_path.endswith('.json'):
                self.config = json.load(f)
            else:
                raise ValueError(f"Unsupported config format: {config_path}")
        
        self.environments: Dict[str, EnvironmentConfig] = {}
        self._parse_config()
    
    def _parse_config(self):
        """Parse configuration file"""
        for env_name, env_config in self.config.get('environments', {}).items():
            services = []
            for svc_config in env_config.get('services', []):
                services.append(ServiceConfig(
                    name=svc_config['name'],
                    type=svc_config['type'],
                    path=svc_config['path'],
                    build_command=svc_config.get('build_command'),
                    start_command=svc_config.get('start_command'),
                    health_endpoint=svc_config.get('health_endpoint'),
                    health_timeout=svc_config.get('health_timeout', 30),
                    env_vars=svc_config.get('env_vars', {}),
                    port=svc_config.get('port'),
                    replicas=svc_config.get('replicas', 1)
                ))
            
            self.environments[env_name] = EnvironmentConfig(
                name=env_name,
                hosts=env_config.get('hosts', []),
                ssh_user=env_config['ssh_user'],
                ssh_key=env_config['ssh_key'],
                deploy_path=env_config['deploy_path'],
                backup_path=env_config.get('backup_path', '/tmp/backups'),
                services=services,
                env_vars=env_config.get('env_vars', {})
            )

# ============================================================================
# Logger
# ============================================================================

class DeploymentLogger:
    """Custom logger with color support"""
    
    COLORS = {
        'DEBUG': '\033[36m',      # Cyan
        'INFO': '\033[32m',       # Green
        'WARNING': '\033[33m',    # Yellow
        'ERROR': '\033[31m',      # Red
        'CRITICAL': '\033[35m',   # Magenta
        'RESET': '\033[0m'
    }
    
    def __init__(self, name: str, log_file: Optional[str] = None):
        self.logger = logging.getLogger(name)
        self.logger.setLevel(logging.DEBUG)
        
        # Console handler with colors
        console_handler = logging.StreamHandler()
        console_handler.setLevel(logging.DEBUG)
        console_handler.setFormatter(self._ColoredFormatter())
        self.logger.addHandler(console_handler)
        
        # File handler
        if log_file:
            file_handler = logging.FileHandler(log_file)
            file_handler.setLevel(logging.DEBUG)
            file_handler.setFormatter(logging.Formatter(
                '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
            ))
            self.logger.addHandler(file_handler)
    
    class _ColoredFormatter(logging.Formatter):
        def format(self, record):
            levelname = record.levelname
            if levelname in DeploymentLogger.COLORS:
                color = DeploymentLogger.COLORS[levelname]
                reset = DeploymentLogger.COLORS['RESET']
                record.levelname = f"{color}{levelname}{reset}"
            return super().format(record)
    
    def debug(self, msg, *args, **kwargs):
        self.logger.debug(msg, *args, **kwargs)
    
    def info(self, msg, *args, **kwargs):
        self.logger.info(msg, *args, **kwargs)
    
    def warning(self, msg, *args, **kwargs):
        self.logger.warning(msg, *args, **kwargs)
    
    def error(self, msg, *args, **kwargs):
        self.logger.error(msg, *args, **kwargs)
    
    def critical(self, msg, *args, **kwargs):
        self.logger.critical(msg, *args, **kwargs)

# ============================================================================
# SSH Connection Manager
# ============================================================================

class SSHManager:
    """Manage SSH connections to remote hosts"""
    
    def __init__(self, logger: DeploymentLogger):
        self.logger = logger
        self.connections: Dict[str, paramiko.SSHClient] = {}
    
    def connect(self, host: str, user: str, key_path: str) -> bool:
        """Establish SSH connection"""
        try:
            client = paramiko.SSHClient()
            client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
            client.connect(
                host,
                username=user,
                key_filename=key_path,
                timeout=30
            )
            self.connections[host] = client
            self.logger.info(f"Connected to {user}@{host}")
            return True
        except Exception as e:
            self.logger.error(f"Failed to connect to {host}: {e}")
            return False
    
    def execute(self, host: str, command: str, timeout: int = 60) -> Tuple[int, str, str]:
        """Execute command on remote host"""
        if host not in self.connections:
            raise Exception(f"No connection to {host}")
        
        client = self.connections[host]
        try:
            stdin, stdout, stderr = client.exec_command(command, timeout=timeout)
            exit_code = stdout.channel.recv_exit_status()
            output = stdout.read().decode('utf-8')
            error = stderr.read().decode('utf-8')
            return exit_code, output, error
        except Exception as e:
            self.logger.error(f"Command execution failed on {host}: {e}")
            return -1, "", str(e)
    
    def upload(self, host: str, local_path: str, remote_path: str) -> bool:
        """Upload file to remote host"""
        if host not in self.connections:
            raise Exception(f"No connection to {host}")
        
        try:
            sftp = self.connections[host].open_sftp()
            sftp.put(local_path, remote_path)
            sftp.close()
            self.logger.debug(f"Uploaded {local_path} -> {host}:{remote_path}")
            return True
        except Exception as e:
            self.logger.error(f"Upload failed to {host}: {e}")
            return False
    
    def download(self, host: str, remote_path: str, local_path: str) -> bool:
        """Download file from remote host"""
        if host not in self.connections:
            raise Exception(f"No connection to {host}")
        
        try:
            sftp = self.connections[host].open_sftp()
            sftp.get(remote_path, local_path)
            sftp.close()
            self.logger.debug(f"Downloaded {host}:{remote_path} -> {local_path}")
            return True
        except Exception as e:
            self.logger.error(f"Download failed from {host}: {e}")
            return False
    
    def close_all(self):
        """Close all connections"""
        for host, client in self.connections.items():
            client.close()
            self.logger.debug(f"Closed connection to {host}")

# ============================================================================
# Build Manager
# ============================================================================

class BuildManager:
    """Manage building of services"""
    
    def __init__(self, logger: DeploymentLogger):
        self.logger = logger
        self.build_dir = Path('/tmp/deploy_builds')
        self.build_dir.mkdir(exist_ok=True)
    
    def build_service(self, service: ServiceConfig, env: EnvironmentConfig) -> Optional[Path]:
        """Build a service and return path to build artifacts"""
        self.logger.info(f"Building {service.name} ({service.type})")
        
        service_path = Path(service.path)
        if not service_path.exists():
            self.logger.error(f"Service path not found: {service.path}")
            return None
        
        build_id = hashlib.md5(f"{service.name}-{datetime.now().isoformat()}".encode()).hexdigest()[:8]
        build_path = self.build_dir / f"{service.name}-{build_id}"
        build_path.mkdir(exist_ok=True)
        
        try:
            if service.type == 'go':
                return self._build_go(service, service_path, build_path)
            elif service.type == 'nodejs':
                return self._build_nodejs(service, service_path, build_path)
            elif service.type == 'python':
                return self._build_python(service, service_path, build_path)
            else:
                self.logger.error(f"Unsupported service type: {service.type}")
                return None
        except Exception as e:
            self.logger.error(f"Build failed for {service.name}: {e}")
            return None
    
    def _build_go(self, service: ServiceConfig, source: Path, dest: Path) -> Path:
        """Build Go service"""
        binary_name = service.name
        binary_path = dest / binary_name
        
        # Build for multiple platforms if needed
        platforms = [
            ('linux', 'amd64'),
            ('linux', 'arm64'),
        ]
        
        for goos, goarch in platforms:
            env = os.environ.copy()
            env.update(service.env_vars or {})
            env['GOOS'] = goos
            env['GOARCH'] = goarch
            
            binary = dest / f"{binary_name}-{goos}-{goarch}"
            if goos == 'windows':
                binary = dest / f"{binary_name}-{goos}-{goarch}.exe"
            
            cmd = ['go', 'build', '-o', str(binary), '-ldflags', '-s -w', '.']
            
            self.logger.debug(f"Running: {' '.join(cmd)}")
            result = subprocess.run(
                cmd,
                cwd=source,
                env=env,
                capture_output=True,
                text=True
            )
            
            if result.returncode != 0:
                self.logger.error(f"Go build failed: {result.stderr}")
                raise Exception(f"Build failed: {result.stderr}")
            
            self.logger.debug(f"Built {binary}")
        
        # Create deployment package
        import tarfile
        package_path = dest.parent / f"{service.name}.tar.gz"
        with tarfile.open(package_path, 'w:gz') as tar:
            for file in dest.iterdir():
                tar.add(file, arcname=file.name)
        
        self.logger.info(f"Built {service.name} successfully")
        return package_path
    
    def _build_nodejs(self, service: ServiceConfig, source: Path, dest: Path) -> Path:
        """Build Node.js service"""
        # Install dependencies
        self.logger.debug("Installing Node.js dependencies")
        result = subprocess.run(
            ['npm', 'ci', '--production'],
            cwd=source,
            capture_output=True,
            text=True
        )
        
        if result.returncode != 0:
            self.logger.error(f"npm install failed: {result.stderr}")
            raise Exception(f"npm install failed: {result.stderr}")
        
        # Run build script if exists
        if service.build_command:
            self.logger.debug(f"Running build command: {service.build_command}")
            result = subprocess.run(
                service.build_command,
                cwd=source,
                shell=True,
                capture_output=True,
                text=True
            )
            
            if result.returncode != 0:
                self.logger.error(f"Build command failed: {result.stderr}")
                raise Exception(f"Build command failed: {result.stderr}")
        
        # Copy to build directory
        import shutil
        shutil.copytree(source, dest / 'source', ignore=shutil.ignore_patterns(
            'node_modules', '.git', 'test', 'tests', '__pycache__', '*.pyc'
        ))
        
        # Create deployment package
        import tarfile
        package_path = dest.parent / f"{service.name}.tar.gz"
        with tarfile.open(package_path, 'w:gz') as tar:
            tar.add(dest, arcname=service.name)
        
        self.logger.info(f"Built {service.name} successfully")
        return package_path
    
    def _build_python(self, service: ServiceConfig, source: Path, dest: Path) -> Path:
        """Build Python service"""
        # Create virtual environment
        self.logger.debug("Creating virtual environment")
        venv_path = dest / 'venv'
        result = subprocess.run(
            ['python3', '-m', 'venv', str(venv_path)],
            capture_output=True,
            text=True
        )
        
        if result.returncode != 0:
            self.logger.error(f"Failed to create virtual environment: {result.stderr}")
            raise Exception(f"Failed to create virtual environment: {result.stderr}")
        
        # Install dependencies
        pip_path = venv_path / 'bin' / 'pip'
        result = subprocess.run(
            [str(pip_path), 'install', '-r', str(source / 'requirements.txt')],
            capture_output=True,
            text=True
        )
        
        if result.returncode != 0:
            self.logger.error(f"pip install failed: {result.stderr}")
            raise Exception(f"pip install failed: {result.stderr}")
        
        # Copy source
        import shutil
        for item in source.iterdir():
            if item.name not in ['venv', '__pycache__', '*.pyc']:
                if item.is_file():
                    shutil.copy2(item, dest)
                else:
                    shutil.copytree(item, dest / item.name)
        
        # Create deployment package
        import tarfile
        package_path = dest.parent / f"{service.name}.tar.gz"
        with tarfile.open(package_path, 'w:gz') as tar:
            tar.add(dest, arcname=service.name)
        
        self.logger.info(f"Built {service.name} successfully")
        return package_path
    
    def cleanup(self):
        """Clean up build directory"""
        import shutil
        shutil.rmtree(self.build_dir, ignore_errors=True)
        self.logger.debug("Cleaned up build directory")

# ============================================================================
# Deployment Manager
# ============================================================================

class DeploymentManager:
    """Main deployment manager"""
    
    def __init__(self, config: DeploymentConfig, logger: DeploymentLogger):
        self.config = config
        self.logger = logger
        self.ssh = SSHManager(logger)
        self.build_manager = BuildManager(logger)
        self.deployment_id = datetime.now().strftime('%Y%m%d-%H%M%S')
        self.deployment_dir = Path('/tmp/deployment')
        self.deployment_dir.mkdir(exist_ok=True)
    
    def deploy(self, environment: str, services: Optional[List[str]] = None, strategy: str = 'rolling'):
        """Deploy to specified environment"""
        if environment not in self.config.environments:
            self.logger.error(f"Unknown environment: {environment}")
            return False
        
        env = self.config.environments[environment]
        self.logger.info(f"Starting deployment to {environment} (strategy: {strategy})")
        
        # Connect to all hosts
        for host in env.hosts:
            if not self.ssh.connect(host, env.ssh_user, env.ssh_key):
                self.logger.error(f"Failed to connect to {host}")
                return False
        
        try:
            # Filter services to deploy
            services_to_deploy = []
            if services:
                service_names = set(services)
                services_to_deploy = [s for s in env.services if s.name in service_names]
                missing = service_names - {s.name for s in services_to_deploy}
                if missing:
                    self.logger.warning(f"Unknown services: {missing}")
            else:
                services_to_deploy = env.services
            
            # Build services
            builds = {}
            for service in services_to_deploy:
                self.logger.info(f"Building {service.name}")
                build_path = self.build_manager.build_service(service, env)
                if not build_path:
                    self.logger.error(f"Failed to build {service.name}")
                    return False
                builds[service.name] = build_path
            
            # Deploy based on strategy
            if strategy == 'blue-green':
                success = self._blue_green_deploy(env, services_to_deploy, builds)
            elif strategy == 'rolling':
                success = self._rolling_deploy(env, services_to_deploy, builds)
            elif strategy == 'canary':
                success = self._canary_deploy(env, services_to_deploy, builds)
            else:
                self.logger.error(f"Unknown deployment strategy: {strategy}")
                return False
            
            if success:
                self.logger.info(f"✅ Deployment to {environment} completed successfully")
                self._notify_success(environment, services_to_deploy)
            else:
                self.logger.error(f"❌ Deployment to {environment} failed")
                self._notify_failure(environment, services_to_deploy)
            
            return success
            
        finally:
            self.ssh.close_all()
            self.build_manager.cleanup()
    
    def _rolling_deploy(self, env: EnvironmentConfig, services: List[ServiceConfig], builds: Dict) -> bool:
        """Perform rolling deployment"""
        self.logger.info("Starting rolling deployment")
        
        for service in services:
            self.logger.info(f"Deploying {service.name}")
            
            for host in env.hosts:
                self.logger.info(f"Deploying to {host}")
                
                # Create backup
                backup_name = f"{service.name}-{self.deployment_id}"
                backup_cmd = f"cp -r {env.deploy_path}/{service.name} {env.backup_path}/{backup_name} 2>/dev/null || true"
                self.ssh.execute(host, backup_cmd)
                
                # Upload new version
                remote_path = f"{env.deploy_path}/{service.name}.tar.gz"
                if not self.ssh.upload(host, str(builds[service.name]), remote_path):
                    self.logger.error(f"Failed to upload {service.name} to {host}")
                    return False
                
                # Extract
                extract_cmd = f"cd {env.deploy_path} && tar -xzf {service.name}.tar.gz && rm {service.name}.tar.gz"
                self.ssh.execute(host, extract_cmd)
                
                # Stop service
                self.ssh.execute(host, f"systemctl --user stop {service.name} || true")
                
                # Start service
                if service.start_command:
                    start_cmd = f"cd {env.deploy_path}/{service.name} && {service.start_command} > /dev/null 2>&1 &"
                    self.ssh.execute(host, start_cmd)
                
                # Health check
                if not self._check_health(host, service):
                    self.logger.error(f"Health check failed for {service.name} on {host}")
                    # Rollback
                    self.logger.info(f"Rolling back {service.name} on {host}")
                    rollback_cmd = f"rm -rf {env.deploy_path}/{service.name} && cp -r {env.backup_path}/{backup_name} {env.deploy_path}/{service.name}"
                    self.ssh.execute(host, rollback_cmd)
                    return False
                
                self.logger.info(f"✅ Deployed {service.name} to {host}")
        
        return True
    
    def _blue_green_deploy(self, env: EnvironmentConfig, services: List[ServiceConfig], builds: Dict) -> bool:
        """Perform blue-green deployment"""
        self.logger.info("Starting blue-green deployment")
        
        for host in env.hosts:
            # Determine current color
            current_color = self._get_current_color(host, env)
            new_color = 'green' if current_color == 'blue' else 'blue'
            
            self.logger.info(f"Current: {current_color}, deploying to {new_color}")
            
            # Deploy to new color
            new_path = f"{env.deploy_path}/{new_color}"
            self.ssh.execute(host, f"mkdir -p {new_path}")
            
            for service in services:
                # Upload to new color
                remote_path = f"{new_path}/{service.name}.tar.gz"
                self.ssh.upload(host, str(builds[service.name]), remote_path)
                
                # Extract
                extract_cmd = f"cd {new_path} && tar -xzf {service.name}.tar.gz && rm {service.name}.tar.gz"
                self.ssh.execute(host, extract_cmd)
            
            # Health check new deployment
            if not self._check_all_health(host, new_path, services):
                self.logger.error(f"Health check failed for new deployment on {host}")
                self.ssh.execute(host, f"rm -rf {new_path}")
                return False
            
            # Switch traffic
            switch_cmd = f"rm -f {env.deploy_path}/current && ln -s {new_path} {env.deploy_path}/current"
            self.ssh.execute(host, switch_cmd)
            
            # Restart services
            for service in services:
                if service.start_command:
                    restart_cmd = f"cd {env.deploy_path}/current/{service.name} && {service.start_command}"
                    self.ssh.execute(host, restart_cmd)
            
            self.logger.info(f"✅ Switched to {new_color} on {host}")
            
            # Clean up old deployment
            old_color = current_color
            self.ssh.execute(host, f"rm -rf {env.deploy_path}/{old_color}")
        
        return True
    
    def _canary_deploy(self, env: EnvironmentConfig, services: List[ServiceConfig], builds: Dict) -> bool:
        """Perform canary deployment"""
        self.logger.info("Starting canary deployment")
        
        if len(env.hosts) < 2:
            self.logger.warning("Canary deployment requires at least 2 hosts, falling back to rolling")
            return self._rolling_deploy(env, services, builds)
        
        # Deploy to first host (canary)
        canary_host = env.hosts[0]
        self.logger.info(f"Deploying canary to {canary_host}")
        
        if not self._deploy_to_host(canary_host, env, services, builds):
            self.logger.error("Canary deployment failed")
            return False
        
        # Monitor canary
        self.logger.info("Monitoring canary for 60 seconds")
        time.sleep(60)
        
        # Check canary health
        for service in services:
            if not self._check_health(canary_host, service):
                self.logger.error("Canary health check failed, rolling back")
                self._rollback_host(canary_host, env, services)
                return False
        
        # Deploy to remaining hosts
        self.logger.info("Canary successful, deploying to remaining hosts")
        
        with concurrent.futures.ThreadPoolExecutor(max_workers=5) as executor:
            futures = []
            for host in env.hosts[1:]:
                futures.append(executor.submit(self._deploy_to_host, host, env, services, builds))
            
            for future in concurrent.futures.as_completed(futures):
                if not future.result():
                    return False
        
        return True
    
    def _deploy_to_host(self, host: str, env: EnvironmentConfig, services: List[ServiceConfig], builds: Dict) -> bool:
        """Deploy all services to a single host"""
        for service in services:
            # Upload
            remote_path = f"{env.deploy_path}/{service.name}.tar.gz"
            if not self.ssh.upload(host, str(builds[service.name]), remote_path):
                return False
            
            # Extract
            extract_cmd = f"cd {env.deploy_path} && tar -xzf {service.name}.tar.gz && rm {service.name}.tar.gz"
            self.ssh.execute(host, extract_cmd)
            
            # Start
            if service.start_command:
                start_cmd = f"cd {env.deploy_path}/{service.name} && {service.start_command} > /dev/null 2>&1 &"
                self.ssh.execute(host, start_cmd)
            
            # Health check
            if not self._check_health(host, service):
                return False
        
        self.logger.info(f"✅ Deployed to {host}")
        return True
    
    def _rollback_host(self, host: str, env: EnvironmentConfig, services: List[ServiceConfig]):
        """Rollback deployment on a host"""
        self.logger.info(f"Rolling back {host}")
        for service in services:
            self.ssh.execute(host, f"rm -rf {env.deploy_path}/{service.name}")
    
    def _get_current_color(self, host: str, env: EnvironmentConfig) -> str:
        """Get current deployment color"""
        result = self.ssh.execute(host, f"readlink {env.deploy_path}/current || echo 'blue'")
        current = result[1].strip()
        return os.path.basename(current) if current else 'blue'
    
    def _check_health(self, host: str, service: ServiceConfig) -> bool:
        """Check service health on a host"""
        if not service.health_endpoint:
            return True
        
        url = f"http://{host}:{service.port}{service.health_endpoint}"
        self.logger.debug(f"Checking health: {url}")
        
        for attempt in range(service.health_timeout):
            try:
                response = requests.get(url, timeout=5)
                if response.status_code == 200:
                    self.logger.debug(f"Health check passed for {service.name} on {host}")
                    return True
            except:
                pass
            
            if attempt < service.health_timeout - 1:
                time.sleep(1)
        
        self.logger.warning(f"Health check failed for {service.name} on {host}")
        return False
    
    def _check_all_health(self, host: str, base_path: str, services: List[ServiceConfig]) -> bool:
        """Check health of all services"""
        for service in services:
            if not self._check_health(host, service):
                return False
        return True
    
    def _notify_success(self, environment: str, services: List[ServiceConfig]):
        """Send success notification"""
        message = f"✅ Deployment to {environment} successful!\n"
        message += f"Services: {', '.join(s.name for s in services)}"
        self.logger.info(message)
        # Could integrate with Slack, email, etc.
    
    def _notify_failure(self, environment: str, services: List[ServiceConfig]):
        """Send failure notification"""
        message = f"❌ Deployment to {environment} failed!\n"
        message += f"Services: {', '.join(s.name for s in services)}"
        self.logger.error(message)

# ============================================================================
# Rollback Manager
# ============================================================================

class RollbackManager:
    """Manage deployment rollbacks"""
    
    def __init__(self, config: DeploymentConfig, logger: DeploymentLogger):
        self.config = config
        self.logger = logger
        self.ssh = SSHManager(logger)
    
    def rollback(self, environment: str, version: Optional[str] = None):
        """Rollback to previous version"""
        if environment not in self.config.environments:
            self.logger.error(f"Unknown environment: {environment}")
            return False
        
        env = self.config.environments[environment]
        self.logger.info(f"Rolling back {environment} to {version or 'previous version'}")
        
        # Connect to hosts
        for host in env.hosts:
            if not self.ssh.connect(host, env.ssh_user, env.ssh_key):
                return False
        
        try:
            if version:
                # Rollback to specific version
                for host in env.hosts:
                    restore_cmd = f"cp -r {env.backup_path}/{version}/* {env.deploy_path}/"
                    self.ssh.execute(host, restore_cmd)
            else:
                # Rollback to latest backup
                for host in env.hosts:
                    # Find latest backup
                    find_cmd = f"ls -t {env.backup_path} | head -1"
                    _, latest, _ = self.ssh.execute(host, find_cmd)
                    latest = latest.strip()
                    
                    if latest:
                        restore_cmd = f"cp -r {env.backup_path}/{latest}/* {env.deploy_path}/"
                        self.ssh.execute(host, restore_cmd)
                        self.logger.info(f"Rolled back {host} to {latest}")
            
            self.logger.info(f"✅ Rollback of {environment} completed")
            return True
            
        finally:
            self.ssh.close_all()

# ============================================================================
# Main CLI
# ============================================================================

def main():
    parser = argparse.ArgumentParser(description='Deployment Script')
    parser.add_argument('--config', '-c', default='deploy-config.yaml',
                        help='Configuration file')
    parser.add_argument('--environment', '-e', required=True,
                        help='Target environment')
    parser.add_argument('--services', '-s', nargs='+',
                        help='Services to deploy (default: all)')
    parser.add_argument('--strategy', choices=['rolling', 'blue-green', 'canary'],
                        default='rolling', help='Deployment strategy')
    parser.add_argument('--action', choices=['deploy', 'rollback', 'status', 'list-backups'],
                        default='deploy', help='Action to perform')
    parser.add_argument('--version', help='Version to rollback to')
    parser.add_argument('--log-file', help='Log file path')
    parser.add_argument('--verbose', '-v', action='store_true',
                        help='Verbose output')
    
    args = parser.parse_args()
    
    # Setup logging
    log_level = logging.DEBUG if args.verbose else logging.INFO
    logger = DeploymentLogger('deploy', args.log_file)
    
    try:
        # Load configuration
        config = DeploymentConfig(args.config)
        logger.info(f"Loaded configuration from {args.config}")
        
        if args.action == 'deploy':
            manager = DeploymentManager(config, logger)
            success = manager.deploy(args.environment, args.services, args.strategy)
            sys.exit(0 if success else 1)
            
        elif args.action == 'rollback':
            manager = RollbackManager(config, logger)
            success = manager.rollback(args.environment, args.version)
            sys.exit(0 if success else 1)
            
        elif args.action == 'status':
            # Show deployment status
            env = config.environments[args.environment]
            logger.info(f"Environment: {args.environment}")
            logger.info(f"Hosts: {', '.join(env.hosts)}")
            logger.info(f"Services: {', '.join(s.name for s in env.services)}")
            
        elif args.action == 'list-backups':
            # List available backups
            env = config.environments[args.environment]
            ssh = SSHManager(logger)
            
            for host in env.hosts:
                if ssh.connect(host, env.ssh_user, env.ssh_key):
                    _, backups, _ = ssh.execute(host, f"ls -l {env.backup_path}")
                    logger.info(f"\nBackups on {host}:")
                    logger.info(backups)
                    ssh.close_all()
        
    except KeyboardInterrupt:
        logger.info("\nDeployment interrupted by user")
        sys.exit(1)
    except Exception as e:
        logger.error(f"Deployment failed: {e}")
        if args.verbose:
            import traceback
            traceback.print_exc()
        sys.exit(1)

if __name__ == '__main__':
    main()