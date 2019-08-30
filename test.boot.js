const assert = chai.assert

var bodyHasClass = function(c) {
	return document.documentElement.classList.contains(c)
}

describe('boot', function() {
	it('should have js classes by default', function() {
	  assert.equal(bodyHasClass("wb-js"), true)
	  assert.equal(bodyHasClass("wb-wasm"), true)
	})

	it('eventually should have run the wasm', function(done){
		var f = setInterval(function(){
			if (bodyHasClass("wb-progress-10") && bodyHasClass("launched") && bodyHasClass("wb-loaded")) {
				clearInterval(f)
				done()
			}
		}, 100)
	})

	// @TODO assert console warn if no main was specified 
	// @TODO assert warn when no content-length is specified 
	// @TODO assert correct fetch config loading
})
  
