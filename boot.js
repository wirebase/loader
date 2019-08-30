document.documentElement.className += " wb-js";
if (typeof WebAssembly === "object") {
  	document.documentElement.className += " wb-wasm";	
  	window.addEventListener("DOMContentLoaded", function(event) {
		var el = document.createElement("script"); el.src = "./latest.js";
		document.body.appendChild(el);
	});
}