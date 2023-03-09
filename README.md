# go-kit generics experiment

From: 

```go
type Endpoints struct {
    PostProfileEndpoint   endpoint.Endpoint
    GetProfileEndpoint    endpoint.Endpoint
    PutProfileEndpoint    endpoint.Endpoint
    PatchProfileEndpoint  endpoint.Endpoint
    DeleteProfileEndpoint endpoint.Endpoint
    GetAddressesEndpoint  endpoint.Endpoint
    GetAddressEndpoint    endpoint.Endpoint
    PostAddressEndpoint   endpoint.Endpoint
    DeleteAddressEndpoint endpoint.Endpoint
}

func MakePostProfileEndpoint(s Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (response interface{}, err error) {
        req := request.(postProfileRequest)
        e := s.PostProfile(ctx, req.Profile)
        return postProfileResponse{Err: e}, nil
    }
}

func decodePostProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
    var req postProfileRequest
    if e := json.NewDecoder(r.Body).Decode(&req.Profile); e != nil {
        return postProfileRequest{}, e
    }
    return req, nil
}

r.Methods("POST").Path("/profiles/").Handler(httptransport.NewServer(
    e.PostProfileEndpoint,
    decodePostProfileRequest,
    encodeResponse,
    options...,	
)
```

To:

```go
type Endpoints struct {
    PostProfileEndpoint   tendpoint.Endpoint[postProfileRequest, postProfileResponse]
    GetProfileEndpoint    tendpoint.Endpoint[getProfileRequest, getProfileResponse]
    PutProfileEndpoint    tendpoint.Endpoint[putProfileRequest, putProfileResponse]
    PatchProfileEndpoint  tendpoint.Endpoint[patchProfileRequest, patchProfileResponse]
    DeleteProfileEndpoint tendpoint.Endpoint[deleteProfileRequest, deleteProfileResponse]
    GetAddressesEndpoint  tendpoint.Endpoint[getAddressesRequest, getAddressesResponse]
    GetAddressEndpoint    tendpoint.Endpoint[getAddressRequest, getAddressResponse]
    PostAddressEndpoint   tendpoint.Endpoint[postAddressRequest, postAddressResponse]
    DeleteAddressEndpoint tendpoint.Endpoint[deleteAddressRequest, deleteAddressResponse]
}

func MakePostProfileEndpoint(s Service) tendpoint.Endpoint[postProfileRequest, postProfileResponse] {
    return func(ctx context.Context, req postProfileRequest) (response postProfileResponse, err error) {
        e := s.PostProfile(ctx, req.Profile)
        return postProfileResponse{Err: e}, nil
    }
}

func decodePostProfileRequest(_ context.Context, r *http.Request) (request postProfileRequest, err error) {
    var req postProfileRequest
    if e := json.NewDecoder(r.Body).Decode(&req.Profile); e != nil {
        return postProfileRequest{}, e
    }
    return req, nil
}

r.Methods("POST").Path("/profiles/").Handler(thttptransport.NewServer(
    e.PostProfileEndpoint,
    decodePostProfileRequest,
    encodeResponse,
    options...,
)
```
