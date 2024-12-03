//Refirma Invoker necesita este div oara funcionar
let componente=document.getElementById("addComponent")
if(componente==null){

    let elemDiv = document.createElement('div');
    elemDiv.setAttribute("id", "addComponent");
    document.body.appendChild(elemDiv);
}


//Firma Peru necesita para no generar conflicto con jquery
const scriptElement = document.createElement('script');
scriptElement.type = 'text/javascript';

scriptElement.textContent = `
    console.log("jQuery demo: " + $.fn.jquery);
    var jqFirmaPeru = jQuery.noConflict(true);
`;
document.head.appendChild(scriptElement);