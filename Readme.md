1. Create the Secret

oc -n test create secret generic ticket-check-proxy --from-literal=session_secret=$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c43)

2. Create Service Account

oc -n test create serviceaccount ticket-check
oc -n test annotate serviceaccount ticket-check serviceaccounts.openshift.io/oauth-redirectreference.ticket-check='{"kind":"OAuthRedirectReference","apiVersion":"v1","reference":{"kind":"Route","name":"ticket-check-authenticated"}}'

3. Create Service

oc -n test apply -f Service.yaml

4. Create Deployment

oc -n test apply -f deployment.yaml

5. Create a route

oc -n test create route reencrypt ticket-check-authenticated --service=ticket-check --port=proxy --insecure-policy=Redirect
