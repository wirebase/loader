document.documentElement.className += " wb-js";
if (typeof WebAssembly === "object") {
  	document.documentElement.className += " wb-wasm";	
  	window.addEventListener("DOMContentLoaded", function(event) {
		el = document.createElement("script"); el.src = "/main.js";
		document.body.appendChild(el);
	});
}