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
                    'primerNombre': 'Victor',
                    'segundoNombre': 'Samuel',
                    'primerApellido': 'Mosquera',
                    'segundoApellido': 'A.',
                    'identificacion': '1087998004',
                    'tipoIdentificacion': 'CC'
                },
                'paciente': {
                    'primerNombre': f_data.get('pnombre'),
                    'segundoNombre': f_data.get('snombre'),
                    'primerApellido': f_data.get('papellido'),
                    'segundoApellido': f_data.get('sapellido'),
                    "identificacion": f_data.get("cedula"),
                    'fechaNacimiento': f_data.get("nacimiento") + "T00:00:00-05:00",
                },
                'receptor': f_data.get('eps'),
                'servicio': f_data.get('servicio'),
                'telefonoPaciente': f_data.get('telefonoPaciente'),
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
