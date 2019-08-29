(function(){
  go = new Go();  
  
  wasm = document.documentElement.dataset.wbMain
  if (!wasm) {
    console.warn("loader: no Wasm file to load specified with the html element's data-wb-main attribute")
    return
  }

  fetch(wasm)
    .then(function(response) {
      if (!response.body) { 
        console.warn('loader: ReadableStream not available on Wasm response, progress will not be shown.')
        return response //no stream to keep progress on
      }

      // if there is a special decompressed-content-length header we can directly use that for our total
      var total = parseInt(response.headers.get('x-decoded-content-length'), 10)
      if (!total) {
        total = parseInt(response.headers.get('content-length'), 10)
        if (!total) {
          console.warn('loader: server did not provide a valid "Content-Length" or "X-Decompressed-Content-Length" header for Wasm binary, progress will not be shown.')
          return response
        }

        // else, if it's 'gzip' or 'brotli' encoded we're gonna estimate it based 
        // on average factors here: https://github.com/golang/go/wiki/WebAssembly#reducing-the-size-of-wasm-files
        var enc = response.headers.get('content-encoding').toLowerCase()
        if (enc.includes("gzip")) {
          total = Math.trunc(total * (4.07)) // between 3.44 / 4.70
        } else if(enc.includes("br")) {
          total = Math.trunc(total * (5.65)) // 6.67 / 4.64
        }
      }

      // start progress reporting
      loaded = 0;
      progress = 0;
      return new Response(
        new ReadableStream({
          start: function(controller) {

            //define a read function that calls itself until the stream is fully read
            var reader = response.body.getReader();
            function read() {
              return reader.read().then(function(res) {
                if (res.done) {
                  controller.close();
                  return;
                }

                // we keep a progress class for easy usage with HTML 
                loaded += res.value.byteLength;
                fraction = loaded/total
                current = Math.trunc(fraction*10)
                if (current > progress) {
                  document.documentElement.classList.remove("wb-progress-"+progress.toString())
                  document.documentElement.classList.add("wb-progress-"+current.toString())
                  progress = current
                }

                // emit a custom event so script can do something more powerfull if they want 
                var event = document.createEvent('Event');
                event.initEvent('wb-progress', true, true);
                event.total = total 
                event.loaded = loaded
                document.documentElement.dispatchEvent(event)

                controller.enqueue(res.value);
                read();
              }).catch(error => {
                console.error("loader: failed to read from Wasm binary response body: "+error)
                controller.error(error)
              })
            }

            return read(); //start reading
          }
        })
      );
    })
    .then(function(response) {
      return response.arrayBuffer()
    })
    .then(function(bytes) {
      document.documentElement.classList.remove("wb-progress-"+progress.toString())
      document.documentElement.classList.add("wb-progress-10")
      return WebAssembly.instantiate(bytes, go.importObject)
    })
    .then(function(result) {
      go.run(result.instance);
      document.documentElement.classList.add("wb-loaded")
    })
    .catch(function(error) {
      console.error("loader: failed to fetch, compile and run Wasm binary at '"+wasm+"': "+error)
    })
})();