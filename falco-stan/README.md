# NATS Streaming output for Sysdig Falco

This sidecar aims to solve this Github [Issue](https://github.com/falcosecurity/kubernetes-response-engine/issues/11). For more details about NATS and NATS Streaming server check [here](https://docs.nats.io/nats-streaming-concepts/intro).

## Configuration

1. Make sure that you have an installed NATS and NATS Streaming instances in the `default namespace` of your cluster: check the [minimal-setup](https://docs.nats.io/nats-on-kubernetes/minimal-setup).
2. Install falco: 

```
helm install --name-template falco --wait -f falco-values.yaml https://github.com/falcosecurity/charts/releases/download/falco-1.5.5/falco-1.5.5.tgz
```
3. Update `falco-nats` sidecar as following (Falco is running as a Daemonset in the default namespace you can edit it with `kubectl edit ds falco`):

```
      - args:
        - /bin/stan-pub
        - -s
        - nats://nats.default:4222
        - -f
        - /tmp/shared-pipe/nats
        - -p
        - $(MY_POD_NAME)
        env:
        - name: MY_POD_NAME
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.name
        image: abirhamzi/falco-stan:dynamic-topic
        imagePullPolicy: Always
        name: falco-nats
```

### Container image

You have this adapter available as a container image. Its name is *abirhamzi/falco-stan:dynamic-topic*.

### Parameters Reference

* -s: Specifies the NATS Streaming server URL where message will be published.  By default
    is: *nats://nats.nats.svc:4222*

* -f: Specifies the named pipe path where Falco publishes its alerts. By default
    is: */tmp/shared-pipe/nats*
* -p: Specifies the NATS Streaming ClienID that should be unique.
