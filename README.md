# loader
Loads and runs a Wirebase app while allowing feedback for unsupported browsers and slow connections.

## TODO 
- [x] setup a build script for combining morphdom, axios, wasm_exec and our own loader code
- [x] setup a release script that can build the javascript release for every Go version
- [x] test jsdelivr as a cdn
- [ ] setup a test page with an example wasm app
- [ ] test on old browsers
- [ ] write documentation
  - [ ] Explain bootstrap base64 encoding
  - [ ] Example HTML with css classes
  - [ ] Custom event
- [ ] Allow fetch call customization; cors etc
- [ ] Add integrity check, maybe based on filename