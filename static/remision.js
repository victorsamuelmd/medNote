(function() {
    "use strict";

    function RFC3339(d) {
        function pad(n) {
            return n < 10 ? '0' + n : n;
        }
        return d.getUTCFullYear() + '-' +
            pad(d.getMonth() + 1) + '-' +
            pad(d.getDate()) + 'T' +
            pad(d.getHours()) + ':' +
            pad(d.getMinutes()) + ':' +
            pad(d.getSeconds()) + '-05:00';
    }

    var fecha = new Date();


    var hacerRemision = function() {
        var remision = function() {

            var f_data = new FormData(document.querySelector('#remisionForm'));
            var data = {
                'medico': {
                    'primer-nombre': 'Victor',
                    'segundo-nombre': 'Samuel',
                    'primer-apellido': 'Mosquera',
                    'segundo-apellido': 'A.',
                    'identificacion': '1087998004',
                    'tipo-identificacion': 'CC'
                },
                'paciente': {
                    'primer-nombre': f_data.get('pnombre'),
                    'segundo-nombre': f_data.get('snombre'),
                    'primer-apellido': f_data.get('papellido'),
                    'segundo-apellido': f_data.get('sapellido'),
                    "identificacion": f_data.get("cedula"),
                    'fecha-nacimiento': f_data.get("nacimiento") + "T00:00:00-05:00",
                },
                'receptor': f_data.get('eps'),
                'servicio': f_data.get('servicio'),
                'fecha': RFC3339(fecha),
                'contenido': f_data.get('contenido')
            };
            return data;
        };
        var rem = new XMLHttpRequest();
        rem.responseType = "arrayBuffer";
        rem.onload = function(data) {
            var pdf = "data:application/pdf;base64," + rem.response;
            document.querySelector("#some").src = pdf;
        };
        rem.open("post", "/remision");
        rem.send(JSON.stringify(remision()));
    };

    document.querySelector('#remitir').addEventListener('mousedown', hacerRemision);
})();
