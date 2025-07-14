![Go](https://img.shields.io/badge/Golang-1.19-blue.svg?logo=go&longCache=true&style=flat)
# Firma Perú Invoker Integration - SGTD PCM
Implementación del Motor de Firma Digital - Firma Perú Invoker Integration - de la [SGTD PCM](https://www.gob.pe/22273)

***Firma Perú Invoker*** es gratuito para las Entidades Públicas del Perú.

Firma Perú Invoker se usa unicamente con certificado digital de persona juridica que se entregan a trabajadores del sector público, certificado digital de persona natural con DNIe o certificado digital emitido por una entidad privada para personas.

La versión [v1.5.1](https://github.com/jumanor/firmaperu-invoker/tree/v1.5.1) es el último lanzamiento

Esta implementación es muy similar a **Refirma Invoker** por lo que los tutoriales de este puede servir aún de guía en **Firma Perú Invoker**

Para mayor información de esta implementación puede ver los siguiente videos: [video 01](https://www.youtube.com/watch?v=bl4OZGWS0lk),[video 02](https://www.youtube.com/watch?v=aOto5CStZNA),[video 03](https://youtu.be/pLlIKqFY8eE)

# Características 
- Soporte para firmar varios documentos 
- Api Rest, puede integrarse en cualquier proyecto web (Php, Python, Java, etc)
- Json Web Tokens (JWT)
- Soporte para protocolo https (SSL/TLS)

# Documentos de la Implementación
- Documento de Integración Firma Peru http://bit.ly/49Ardlv
- Guía para el uso e integración de la Plataforma Nacional de Firma Digital en la Administración Pública https://cutt.ly/urmqJxHM

# Cómo solicitar identificadores? 

Para ejecutar *Firma Perú Invoker Integration* es necesario que el Responsable de la Oficina de Tecnologias solicite
las credenciales de acceso a la *Secretaria de Gobierno y Transformación Digital* de la *PCM*

La solicitud la puede realizar en https://www.gob.pe/22273 y busca la opción: *"Formato de solicitud de credenciales de Firma Perú - Firmador web"*

Se le proporicionara un archivo **fwAuthorization.json** dentro del cual se encuentran las credenciales
**[clientId]** y **[clientSecret]** para el uso de Firma Peru Invoker en el Sistema de Gestión Documental de su institución.   

# Probando con Docker

1) Levantamos un contenedor de refirma-invoker
```
docker run -d --name firmaperu-invoker -p 5050:5050 -p 9091:9091 -e CLIENT_ID=mi_client_id -e CLIENT_SECRET=mi_cliente_secret jumanor/firmaperu-invoker:v1.5.1
```
2) Probamos el **example01** (el proceso de firma de los clientes solo esta disponible para Sistema Operativo Windows)
```
http://127.0.0.1:5050/example01/test.html
```

# Instalación del Servidor

Esta disponible un video de la instalación en el siguiente [enlace](https://www.youtube.com/watch?v=7q4dS8y3Sws)

### Requisito
Esta implementación usa [7-zip](https://www.7-zip.org/) que normalmente ya viene instalada en **LINUX**; sin embargo, en **WINDOWS** tendra que instalar manualmente y verificar que se puede acceder desde el terminal.

Windows
    
    C:\Users\Jumanor>7z i

Linux

    jumanor@ubuntu:~$7z i

Para instalar 7z en **Centos 8** seguir los siguientes pasos:
1) Abrir un terminal
2) yum install epel-release         *(instalamos repositorio epel)*
3) yum install p7zip                *(instalamos p7zip)*
4) ln -s /usr/bin/7za /usr/bin/7z   *(creamos enlace simbolico a 7z)*

Para instalar 7z en **Ubuntu/Debian** seguir los siguientes pasos:
1) Abrir un terminal
2) apt install update              
3) apt install p7zip-full          *(instalamos p7zip)*

Para instalar 7z en **Windows 10/11** seguir los siguientes pasos:
1) Descargar e instalar 7z de [aquí](https://www.7-zip.org/)
2) La ruta de instalacion por defecto es C:\Program Files\7-Zip
4) Abrir un cmd *(simbolo de sistema o consola de comandos)*
3) setx path "%path%;C:\Program Files\7-Zip"    *(actualizamos la variable de entorno path)*

### Instalación
Se compilo *Firma Perú Invoker Integration* para Windows y Linux, y estan disponibles en los [releases](https://github.com/jumanor/firmaperu-invoker/releases/latest).

1. Descargue el ejecutable
   
   Windows 64-bit: [main.exe](https://github.com/jumanor/firmaperu-invoker/releases/latest/download/main.exe)
   
   Linux 64-bit:   [main](https://github.com/jumanor/firmaperu-invoker/releases/latest/download/main)

2. Copia la carpeta **public** del repositorio esta contiene 2 imagenes: iFirma.png e iLogo.png
3. Crea un archivo **config.properties** con los siguientes parametros :
    ``` bash
    # Identificador proporcionado por SEGDI-PCM
    clientId=K57845459hkj
    # Identificador proporcionado por SEGDI-PCM
    clientSecret=TYUOPDLDFDG
    # Direccion Ip y Puerto de escucha Firma Perú Invoker Integration
    serverAddress=0.0.0.0:9091
    # Clave secreta para generar Tokens
    secretKeyJwt=muysecretokenjwt
    # Usuario que accedera a la API
    userAccessApi=usuarioAccesoApi
    # Tiempo de expiración del Token en minutos. Ejemplo 5 minutos (Opcional)
    timeExpireToken=5
    # Certificado SSL/TLS (Opcional)
    #certificateFileTls=C:\Users\jumanor\cert.pem
    #certificateFileTls=/home/jumanor/cert.pem
    # Clave Privada SSL/TLS  (Opcional)
    #privateKeyFileTls=C:\Users\jumanor\key.pem
    #privateKeyFileTls=/home/jumanor/key.pem
    ``` 
4. En caso desee habilitar protocolo **https** es necesario que ingrese los siguientes parametros :
    ``` bash
    # Certificado SSL/TLS (Opcional)
    certificateFileTls=/etc/letsencrypt/live/midominio.com/fullchain.pem
    # Clave Privada SSL/TLS (Opcional)
    privateKeyFileTls=/etc/letsencrypt/live/midominio.com/privkey.pem
    ```
5. Si necesitas controlar la rotación del log utiliza los siquientes parametros:
    ``` bash
    # Tamaño máximo en MB antes de rotar (Opcional)
    maxSize=10
    # Número máximo de archivos de backup (Opcional)
    maxBackups=3
    # Máximo de días a mantener los logs (Opcional)
    maxAge=3
    ```
6. Ejecuta *ReFirma Invoker Integration*

    Windows

        main.exe

    Linux

        ./main

# Instalación del Cliente
Estan disponibles videos del funcionamiento (ejemplos) en los siguientes enlaces: [enlace 01](https://www.youtube.com/watch?v=GPdfa7NeKZw).

Firma Perú Invoker usa **Microsoft Click Once** para invocar a Firma Perú (Componente Web).

1. Si esta usando navegador **Chome** o **Firefox** instala los siguientes plugins para habilitar **Microsoft Click Once**:

    - Chrome instale este [plugin](https://chromewebstore.google.com/detail/cegid-peoplenet-clickonce/jkncabbipkgbconhaajbapbhokpbgkdc) 

    - Firefox instale este [plugin](https://addons.mozilla.org/es/firefox/addon/breez-clickonce)  
    
2. En caso use el navegador **Edge** no es necesario instalar nada adicional (Recomendable).

3. Copia la carpeta [example](https://github.com/jumanor/firmaperu-invoker/tree/master/example) de este repositorio en un Servidor Web (ver el siguiente [video](https://youtu.be/7q4dS8y3Sws?t=218) para mayor detalle)

    3.1. En caso use **Visual Studio Code** instale el plugin [Live Server](https://marketplace.visualstudio.com/items?itemName=ritwickdey.LiveServer) que habilita un Servidor Web Embebido (Recomendable).

5. Ingresa a cualquier ejemplo que desee probar ejecutando **`http://direccion_ip:puerto/example01/test.html`**

Nota

Las siguientes instalaciones de Fima Perú (Componente Web) no son Oficiales de la [SGTD PCM](https://www.gob.pe/22273)

- Si usa **Chome** o **Firefox** es posible ejecutar Firma Perú (Componente Web) sin instalar plugins(extensiones) adicionales para detalles de la instalación ver el siguiente [video](https://www.youtube.com/watch?v=3krIhVr6NCs).
- Tambien es posible ejecutar Firma Perú (Componente Web) en cualquier derivado de **Ubuntu (Linux)** para detalles de la instalación ver el siguiente [video](https://www.youtube.com/watch?v=-YgnULCkjlk). 

# Funcionamiento

A continuación un manera simplicada del uso de Firma Perú Invoker Integration con **JavaScript** del lado del Cliente:


``` javascript
//Listamos los documentos que se desean firmar digitalmente
let pdfs=[];
pdfs[0]={url:"http://miservidor.com/docs1.pdf",name:"doc1"};
pdfs[1]={url:"http://miservidor.com/docs2.pdf",name:"doc2"};

//Enviamos la posicion en donde se ubicara la representación gráfica de la firma digital
let firmaParam={};
firmaParam.posx=10;
firmaParam.posy=12;
firmaParam.reason="Soy el autor del documento pdf";
firmaParam.role="Programador Full Stack";
firmaParam.stampSigned="http://miservidor.com/estampillafirma.png";//parametro opcional
firmaParam.pageNumber=1; //parametro opcional, pagina donde se pondra la firma visible 
firmaParam.visiblePosition=false;//parametro opcional, interfaz gráfica(posicion de firma) nativo de firma perú
firmaParam.oneByOne=false;//parametro opcional,false: una sola vez ubica la posicion de la firma true: por cada documento ubica la posicion de la firma solo si visiblePosition=true
firmaParam.signatureStyle=1;//parametro opcional,0:sin representacion grafica 1:horizontal 2:vertical 3:solo estampado 4:solo descripción
firmaParam.stampTextSize=14;//parametro opcional,
firmaParam.stampWordWrap=37;//parametro opcional,

//Llamamos a Firma Perú Invoker Integration con la dirección ip en donde se ejecuta main.exe o main
let firma=new FirmaPeru("http://192.168.1.10:9091");
//Muy importante !!!
//El Sistema de Gestion Documental se encarga de la autenticación y envía un token al Cliente
//Este método se usa solo como demostración no se debe de usar en el Cliente
let token=await firma.autenticacion("usuarioAccesoApi");
//Realiza el proceso de Firma Digital
let url_base=await firma.ejecutar(pdfs,firmaParam,token);

//En este caso obtenemos los documentos firmados digitalmente y los enviamos a un frame
document.getElementById("frame1").src=url_base+"/"+encodeURI("doc1")+"/"+encodeURI(token);
document.getElementById("frame2").src=url_base+"/"+encodeURI("doc2")+"/"+encodeURI(token);
```          

# Integrando a un Sistema de Gestion Documental

Esta implementación de *Firma Perú Invoker Integration* se puede usar en ***cualquier proyecto web*** (Php, Java, Python, etc) solo tiene que consumir las Api Rest implementadas, para controlar el acceso se usa JSON Web Tokens ([JWT](https://jwt.io/))

El *Sistema de Gestión Documental* autentica a los Usuarios normalmente contra una Base de Datos,
despues de la autencación satisfactoria se debe de consumir  el API REST **/autenticacion** de Firma Perú Invoker 
y enviar el **token** al Cliente.

![a link](https://raw.githubusercontent.com/jumanor/firmaperu-invoker/master/public/funcionamiento.jpeg)

A continuacion algunos ejemplos de captura del **token de autenticación** en el lado del servidor:

Ejemplo con Curl
``` bash
 curl -X POST http://127.0.0.1:9091/autenticacion -H "Content-Type: application/json; charset=UTF-8" -d '{"usuarioAccesoApi":"usuarioAccesoApi"}'
```

Ejemplo con Python
``` python
import requests
import json
api_url = "http://127.0.0.1:9091/autenticacion"
param={"usuarioAccesoApi":"usuarioAccesoApi"}
response = requests.post(api_url,json=param)
if response.status_code == 200:
	token=response.json().get("data")
	print(token)

```
Ejemplo con Php
``` php
$params=array("usuarioAccesoApi"=>"usuarioAccesoApi");
$postdata=json_encode($params);
$opts = array('http' =>
    array(
    'method' => 'POST',
    'header' => 'Content-type: application/json',
    'content' => $postdata
    )
);
$context = stream_context_create($opts);
@$response = file_get_contents("http://127.0.0.1:9091/autenticacion", false, $context);
if(isset($http_response_header) == true){
    
    $status_line = $http_response_header[0];
    preg_match('{HTTP\/\S*\s(\d{3})}', $status_line, $match);
    $status = $match[1];

    if ($status == 200){
        $obj=json_decode($response,true);
        $token=$obj["data"];
        echo $token;
    }    
}
```

# Contribución

Por favor contribuya usando [Github Flow](https://guides.github.com/introduction/flow/). Crea un fork, agrega los commits, y luego abre un [pull request](https://github.com/fraction/readme-boilerplate/compare/).

# License
Copyright © 2024 [Jorge Cotrado](https://github.com/jumanor). <br />
This project is [MIT](https://github.com/jumanor/refirmainvoker/blob/master/License) licensed.
