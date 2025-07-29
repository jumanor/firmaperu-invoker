%define name firmaperu-invoker
%define version 1.5.1
%define release 1
%define buildroot %{_topdir}/BUILDROOT
%define _binary_payload w9.gzdio

Name: %{name}
Version: %{version}
Release: %{release}%{?dist}
Summary: Motor de Firma Digital - Firma Perú Invoker Integration
Group: Applications/Internet
License: MIT
URL: https://github.com/jumanor/firmaperu-invoker
Requires: p7zip

BuildRoot: %{buildroot}

%description
Firma Perú Invoker es una implementación del Motor de Firma Digital para la Secretaria de Gobierno y Transformación Digital - PCM.
Este servicio permite la integración con Sistemas de Gestión Documental para la firma digital de documentos PDF.

%prep
echo "Preparando el paquete..."

%build
echo "No se requiere compilación, usando binarios pre-compilados..."

%install
mkdir -p %{buildroot}/opt/firmaperu-invoker
mkdir -p %{buildroot}/etc/systemd/system
mkdir -p %{buildroot}/var/log/firmaperu-invoker
install -m 755 %{_sourcedir}/main %{buildroot}/opt/firmaperu-invoker/
install -m 644 %{_sourcedir}/config.properties.example %{buildroot}/opt/firmaperu-invoker/
cp -r %{_sourcedir}/public %{buildroot}/opt/firmaperu-invoker/
install -m 644 %{_sourcedir}/firmaperu-invoker.service %{buildroot}/etc/systemd/system/

%pre
# Crear el usuario firmaperu si no existe
if ! getent passwd firmaperu >/dev/null; then
    useradd -r -s /sbin/nologin -d /opt/firmaperu-invoker firmaperu
fi

%post
if [ $1 -eq 1 ] ; then
    # Primera instalación
    echo "Configurando servicio firmaperu-invoker..."
    
    # Copiar el archivo de configuración de ejemplo
    if [ ! -f /opt/firmaperu-invoker/config.properties ]; then
        cp /opt/firmaperu-invoker/config.properties.example /opt/firmaperu-invoker/config.properties
        echo "Se ha creado el archivo de configuración a partir del ejemplo"
    fi
    
    # Asignar permisos adecuados
    chown -R firmaperu:firmaperu /opt/firmaperu-invoker
    chown -R firmaperu:firmaperu /var/log/firmaperu-invoker
    
    # Crear enlace simbólico si es necesario
    if [ ! -f /usr/bin/7z ] && [ -f /usr/bin/7za ]; then
        ln -s /usr/bin/7za /usr/bin/7z
    fi

    systemctl daemon-reload
    systemctl enable firmaperu-invoker.service
    echo "El servicio ha sido habilitado pero no iniciado."
    echo "Para iniciar el servicio ejecute: systemctl start firmaperu-invoker"
    echo "Para configurar el servicio, edite: /opt/firmaperu-invoker/config.properties"
fi

%preun
if [ $1 -eq 0 ] ; then
    # Desinstalación
    systemctl stop firmaperu-invoker.service
    systemctl disable firmaperu-invoker.service
fi

%files
%defattr(-,root,root,-)
%dir /opt/firmaperu-invoker
/opt/firmaperu-invoker/main
%config(noreplace) /opt/firmaperu-invoker/config.properties.example
/opt/firmaperu-invoker/public
/etc/systemd/system/firmaperu-invoker.service
%dir /var/log/firmaperu-invoker

%changelog
* Wed Jul 23 2025 Jorge Cotrado <jorge.cotrado@gmail.com> - 1.5.1-1
- Versión inicial del paquete RPM para Firma Perú Invoker
