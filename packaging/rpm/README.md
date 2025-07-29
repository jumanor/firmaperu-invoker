# Instalación de Firma Perú Invoker mediante RPM

Este directorio contiene los archivos necesarios para generar un paquete RPM de Firma Perú Invoker, lo que facilita la instalación en sistemas basados en Red Hat como CentOS, RHEL, Fedora, Rocky Linux, Alma Linux, etc.

## Requisitos previos

Para construir el paquete RPM, necesita tener instaladas las herramientas de construcción de RPM:

```bash
# En CentOS/RHEL/Rocky Linux/Alma Linux
sudo dnf install rpm-build rpmdevtools -y
```

## Construir el paquete RPM

1. Ejecute el script de construcción:

```bash
chmod +x build-rpm.sh
./build-rpm.sh
```

2. El paquete RPM generado se encontrará en el directorio `dist/`.

## Requisitos previos de instalación del RPM
Ante de instalar el paquete RPM, es necesario instalar 7zip

```bash
sudo dnf install epel-release -y
sudo dnf install p7zip -y
```

## Instalar el paquete RPM

Una vez generado el paquete RPM, puede instalarlo con:

```bash
sudo rpm -ivh dist/firmaperu-invoker-1.5.1-1.el7.x86_64.rpm
```

## Configuración post-instalación

Durante la instalación, el archivo de configuración `config.properties.example` se copia automáticamente a `config.properties` si este último no existe.

1. Edite el archivo de configuración según sus necesidades:

```bash
sudo vi /opt/firmaperu-invoker/config.properties
```

2. Edite el archivo para configurar los parámetros requeridos:
   - clientId
   - clientSecret
   - serverAddress
   - secretKeyJwt
   - userAccessApi

3. Inicie el servicio:

```bash
sudo systemctl start firmaperu-invoker
```

4. Para verificar el estado del servicio:

```bash
sudo systemctl status firmaperu-invoker
```

5. Para habilitar el inicio automático del servicio durante el arranque del sistema:

```bash
sudo systemctl enable firmaperu-invoker
```
## Desinstalación

Para desinstalar Firma Perú Invoker:

```bash
sudo rpm -e firmaperu-invoker
```
