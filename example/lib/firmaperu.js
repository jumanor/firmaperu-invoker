export class FirmaPeru{

    ARGUMENTOS=null;
    URL_BASE=null;

    EVENT_SUCCESS_INVOKER=null;
    EVENT_ERROR_INVOKER=null;
    EVENT_INICIO_INVOKER=null;
    //Direccion de Servidor Refirma Invoker
    URL_SERVER_REFIRMA_INVOKER=null;
    
    constructor(url_server_refirma_invoker){
        this.URL_SERVER_REFIRMA_INVOKER=url_server_refirma_invoker;
    }
    event_signatureInit(){
       // dispatchEventClient('sendArguments', this.ARGUMENTOS);
       //   this.EVENT_INICIO_INVOKER
    }
    event_signatureOk(){
        this.EVENT_SUCCESS_INVOKER(this.URL_BASE);
    }
    event_signatureCancel(e){
        console.log(e)
        this.EVENT_ERROR_INVOKER("FIRMA CANCELADA")
    }
    precarga(){
        
        this.event_signatureInit = this.event_signatureInit.bind(this);
        this.event_signatureOk = this.event_signatureOk.bind(this);
        this.event_signatureCancel = this.event_signatureCancel.bind(this);

        window.signatureInit=this.event_signatureInit;
        window.signatureOk=this.event_signatureOk;
        window.signatureCancel=this.event_signatureCancel;


    }////////////////////////////////////////////////////////////////////////////////////
    // Solamente el SGD debe usar /autenticacion.
    // No utilize este metodo en producción
    async autenticacion(usuarioAccesoApi){
        
        let response=await fetch(this.URL_SERVER_REFIRMA_INVOKER+"/autenticacion",{
            method:'POST',
            body:JSON.stringify({usuarioAccesoApi:usuarioAccesoApi}),
            headers: {
                        'Content-Type': 'application/json; charset=UTF-8',      
                    },
        });
        if(!response.ok){//200-299
            console.log(response.statusText)
            let tt=await response.text()
            throw Error(tt);
        }
        
        let result=await response.json();
        return result.data
        
    }////////////////////////////////////////////////////////////////////////////////////
    async ejecutar(urlPdfs,parametros,token){
            this.precarga();
            console.log(window.signatureInit);
        
            let params={};
            //params.pdfs=[];
            //params.pdfs[0]={url:"http://127.0.0.1:5500/01.pdf",name:"doc0"};
            //params.pdfs[1]={url:"http://127.0.0.1:5500/02.pdf",name:"doc1"};
            params.pdfs=urlPdfs
            //params.firma={};
            //params.firma.posx=10;
            //params.firma.posy=12;
            params.firma=parametros;
             
                let response=await fetch(this.URL_SERVER_REFIRMA_INVOKER+"/argumentsServletPCX",{
                    method:'POST',
                    body:JSON.stringify(params),
                    headers: {
                                'Content-Type': 'application/json; charset=UTF-8', 
                                'x-access-token' : token    
                            },
                });

                if(!response.ok){//200-299
                    console.log(response.statusText)
                    let tt=await response.text()
                    throw Error(tt);
                }
                
                let result=await response.json();
                console.log(result);

                let argumentos=result.data.argumentosBase64;
                this.URL_BASE=result.data.urlBase;

                var port="48596";
                startSignature(port,argumentos);
           
            return new Promise((resolve,reject)=>{

                this.EVENT_SUCCESS_INVOKER=resolve;
                this.EVENT_ERROR_INVOKER=reject;

            });

    }//////////////////////////////////////////////////////////////////////////////
}