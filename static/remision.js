f_data = new FormData(document.querySelector('#remisionForm'));

function RFC3339(d){
  function pad(n){
    return n<10 ? '0'+n : n
  }
  return d.getUTCFullYear()+'-'
    + pad(d.getMonth()+1)+'-'
    + pad(d.getDate())+'T'
    + pad(d.getHours())+':'
    + pad(d.getMinutes())+':'
    + pad(d.getSeconds())+'-05:00'
}

fecha = new Date();


var hacerRemision = function(){
  var remision = {
    medico: {
      "primer-nombre": "Victor",
      "segundo-nombre": "Samuel",
      "identificacion": "1087998004",
      "tipo-identificacion": "CC"
    },
    "paciente": {
      "primer-nombre": f_data.get('pnombre')
    },
    "receptor": f_data.get('eps'),
    "servicio": f_data.get('servicio'),
    "fecha": RFC3339(fecha)
  }
  rem = new XMLHttpRequest();
  rem.responseType = "arrayBuffer";
  rem.onload = function(data){
    var pdf = "data:application/pdf;base64," + rem.response;
    document.querySelector("#some").src = pdf;
  }
  rem.open("post", "/remision");
  rem.send(JSON.stringify(remision));
}

document.querySelector('#remitir').addEventListener('mousedown', hacerRemision);
