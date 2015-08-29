// ***************************************************************
//      app.js
// ***************************************************************

var config = require("config");
var asrt = require("./asrt-table");

$(document).ready(function() {
  $('[data-toggle=offcanvas]').click(function() {
    $('.row-offcanvas').toggleClass('active');
  });
});

asrt.Initialize(config.url, config.interval);


