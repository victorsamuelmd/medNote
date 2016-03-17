(function(){
    'use strict';
    var imc, showIMC, verFormula;

    imc = function(peso, talla){
        return peso/(talla*talla);
    }

    showIMC = function(){
        var peso, talla;
        peso = document.querySelector('#peso').value;
        talla = document.querySelector('#talla').value;
        document.querySelector('#imc').value = imc(peso, talla);
    }


    verFormula = function(){
        var data = new FormData(document.getElementsByTagName('form')[0]);
        var reg = new XMLHttpRequest();
        reg.onload = function(){
                var file = window.URL.createObjectURL(reg.response);
                var a = document.createElement("a");
                a.href = file;      
                window.open(file);
                console.log('hola');
        }
        reg.open("POST", '/formula');
        reg.send({
            'nombre': '' + data.get('pnombre') + data.get('snombre') + data.get('papellido'),
            'id': data.get('cedula'),
            'eps': 'Asmetsalud',
            'centro-salud': 'Arboleda',
            'receta': data.get('conducta')
        }); 
    }
    
    document.querySelector('#talla').addEventListener('blur', showIMC);
    document.querySelector('#formula-pdf').addEventListener('click', verFormula);
})();
