# beargrabz

Monitor Kubernetes Clusters for authorisation tokens being passed in clear-text.


### Deploy beargrabz 
```sh
kubectl apply -f deploy.yaml
```

Check `beargrabz` logs to see what bearer tokens it has found so far:
```sh
kubectl logs eavesdropper
```

### Deploy full working example
This will deploy 3 pods:

1. `httpbin` pod and service to serve HTTP on port 80.
2. `api-client` pod to simulate requests.
3. `beargrabz` pod to eavesdrop requests from `api-client` to `httpbin`.

```sh
kubectl apply -f sample/playground.yaml
```

Follow `beargrabz` logs to see requests from `api-client` to `httpbin`:
```sh
kubectl logs eavesdropper -f
```

The result should be a new entry in the log every second:
```
10.244.0.23:34570 -> 10.0.114.142:8000
GET /
Host: httpbin:8000
Authorization: Bearer GoUwqbik432***********
```

### Security Requirements

`beargrabs` has a few security requirements:

- run as root on the container.
- have `NET_ADMIN` capability.
- have access to the host network.

**Running on Kubernetes**
An extract of the yaml configuration is as follows:

```yaml
    securityContext:
      capabilities:
        add: ["NET_ADMIN"]
  hostNetwork: true
```

Check [deploy.yml](deploy.yml) for a working example.

**Running manually with docker**
```sh
docker run --rm --security-opt=no-new-privileges --cap-drop=NET_ADMIN --network="host" paulinhu/beargrabz 
```


## License

Licensed under the MIT License. You may obtain a copy of the License [here](LICENSE).