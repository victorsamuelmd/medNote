(function() {
    'use strict';
    var imc, showIMC, verFormula;
    var currentDate = new Date();
    var formulario = new FormData(document.getElementsByTagName("form")[0]);

    imc = function(peso, talla) {
        return peso / (talla * talla);
    };

    showIMC = function() {
        var peso, talla;
        peso = document.querySelector('#peso').value;
        talla = document.querySelector('#talla').value;
        document.querySelector('#imc').value = imc(peso, talla);
    };

    document.querySelector('#talla').addEventListener('blur', showIMC);
})();
