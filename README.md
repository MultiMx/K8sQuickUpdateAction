# K8sQuickUpdateAction

Quickly update deployments' image action with multi-threaded method.

## Usage

Workloads in same element in "workloads" field will run in parallel. Only when all workloads in one element are deployed will it go to the next element.

Program will throw exit status 1 when any operation failed after all deploy operations completed,

```yaml
      - name: Update Deployments
        uses: MultiMx/K8sQuickUpdateAction@v0.6
        with:
          k8s: |
            prod:
              backend: https://cluster1.com
              token: ${{ secrets.CATTLE_TOKEN_1 }}
            dev:
              backend: https://cluster2.com
              token: ${{ secrets.CATTLE_TOKEN_2 }}
          workloads: |
            -
              an-namespace:
                an-workload:
                  image: ${{ steps.image.outputs.CORE_VERSION }}
                  wait: true
            -
              an-namespace:
                an-workload:
                  image: ${{ steps.image.outputs.AUTH_VERSION }}
                another-workload:
                  image: ${{ steps.image.outputs.FE_VERSION }}
              another-namespace:
                ...
```