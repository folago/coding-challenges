# NOTES

This is what did:

- Finish implementing the `crypto` package by add a hash algorithm, implementing the `Signer()` interface and make a enum for the signing algorithms.
- Implement the in `domain` package the device and the signature types.
- Implement the `persistence` package, define a repository interface and implement it in memory.
- Add a service package that holds the business logic of signing some data and persisting both signed data and devices.
- Implement most of the REST API now that the pieces are available.



## Notes
Few notes collected here.

I focused mostly on the domain and service code so the REST API is a bit neglected. Having time I would have used some OpenAPI code generator. Also there is no mention on response times, timeouts, and input size. The current implementation will not survive a terabyte of data to sign.

I also completely skipped the serialization of the devices, not because it is not important but because I was not useful for a proof of concept with in memory storage. I wanted to have something running in order to be able to think about the system as a whole and to play with it a bit.

The UUID lib was not in the go.mod file and since the `CreateSignatureDevice` had an id in the signature I assumed the id is passed to us and assumed a string is good enough™, at least for the tech test.

### About Q/A Testing

I added a test for the service to simulate the concurrent use of a device. As far as it looks promising I would not take that test as verifying correct behavior.
In order to do that an idea could be to run a lot of concurrent request against the service, via the REST API and store the data somewhere (in a file or a DB). Then read the data and reconstruct the signing history. 
The domain types have all the data available to walk back a trail of signing since the `domain.Signature` holds the signed data which contains the last signature and the signature counter (reminds me of git commits).
This allows us to verify the signatures as well, never mentioned in the specs but probably nice to check.
It seems fun but sadly I had no time to implement it (it also seems a good place for off-by-one errors).


### About Technical Constraints & Considerations

I believe it is mostly taken care of:
- The `service` should be able to run with multiple concurrent clients properly, the test is only partial since it only uses one device, and still a proper verification of correctness is non implemented. A restriction of this implementation is that scaling out is not possible unless multiple service can handle a disjoint sets of devices (maybe consistent hashing?), that being another problem in itself.
- The monotonicity of the signature counter should be satisfied in the case of only one service.
- Adding another signing algorithm should be straightforward with the enum.
- The repositories interface are a good starting point to segregate functionality needed to persist the domain entities, in case different storage are needed for different entities. They also abstract away the actual storage system, switching to RDMS should not be too much hassle.

