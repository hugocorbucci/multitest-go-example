# multitest-go-example

An example of a webapp using the multitest style (I'm looking for a better name, feel free to make suggestions) in golang.

The presentation that explains this work is hosted in https://hugocorbucci.github.io/multitest-go-example/.

If you just want a TL;DR, check [this test](https://github.com/hugocorbucci/multitest-go-example/blob/98471512b728be43bc73454b119002f64dcd73cc/internal/server/http_test.go#L85-L108) which runs in three different forms depending on environment variables:

1) As a pure "unit"/isolated test. In which the test itself just creates a server by instantiating a server with mocked or stubbed dependencies and making direct calls to it.
2) As an "integration" test. In which the test creates a server in a separate go routine and starts it off and then fires regular http request to that local server with out of process (but local) dependencies.
3) As a "smoke" test. In which the test assumes a server is running somewhere and issues requests in a "black box" style and checks the responses only. These can be used to run against "real" environments such as your staging or production environments.

I've found this style of test to be better suited for microservices for both http/json and grpc/protobuf communication protocols. It trades off the obvious split of unit/integration/smokes and the [testing pyramid](https://martinfowler.com/bliki/TestPyramid.html) style for a smaller testing code base (allegedly cheaper to maintain) without necessarily suffering from an increased feedback loop cost (units still run fast, integration midly slower and smokes slower, but you choose which style you want on runtime, not development time).
