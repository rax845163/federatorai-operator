#!/bin/sh

##################################################################
#
#   This script is created for installing Federator.ai Operator
#
#   Usage: ./install.sh $tag_name or ./install.sh
##################################################################

is_pod_ready()
{
  [[ "$(kubectl get po "$1" -n "$2" -o 'jsonpath={.status.conditions[?(@.type=="Ready")].status}')" == 'True' ]]
}

pods_ready()
{
  [[ "$#" == 0 ]] && return 0
  
  namespace="$1"
  
  all_pods="$(kubectl get po -n $namespace -o 'jsonpath={.items[*].metadata.name}')"
  for pod in $all_pods; do
    is_pod_ready $pod $namespace || return 1
  done

  return 0
}

leave_prog()
{
    echo -e "\n$(tput setaf 5)Downloaded YAML files are located under $file_folder $(tput sgr 0)"
    cd $current_location > /dev/null
}

wait_until_pods_ready()
{
  period="$1"
  interval="$2"
  namespace="$3"
  target_pod_number="$4"

  for ((i=0; i<$period; i+=$interval)); do

    if [[ "`kubectl get po -n $namespace 2>/dev/null|wc -l`" -ge "$target_pod_number" ]]; then
        if pods_ready $namespace; then
            echo -e "\nAll $namespace pods are ready."
            return 0
        fi
    fi

    echo "Waiting for $namespace pods to be ready..."
    sleep "$interval"
    
  done

  echo -e "\n$(tput setaf 1)Warning!! Waited for $period seconds, but all pods are not ready yet. Please check $namespace namespace$(tput sgr 0)"
  leave_prog
  exit 4
}

get_grafana_route()
{
    if [[ "$openshift_version" != "" ]] ; then
        link=`oc get route -n $1|grep grafana|awk '{print $2}'`
        echo -e "\n========================================"
        echo "You can now access GUI through $(tput setaf 6)http://${link} $(tput sgr 0)"
        echo "Default login credential is $(tput setaf 6)admin/admin$(tput sgr 0)"
        echo -e "\nAlso, you can start to apply alamedascaler CR for namespace you would like to monitor."
        echo "$(tput setaf 6)Review administration guide for further details.$(tput sgr 0)"
        echo "========================================"
    fi
}

if [ "$1" != "" ];then
    tag_number="$1"
else
    
    while [[ "$tag_correct" != "y" ]]
    do
        read -r -p "$(tput setaf 2)Please input Federator.ai Operator tag:$(tput sgr 0) " tag_number </dev/tty
        default="y"
        read -r -p "$(tput setaf 2)Is tag \"${tag_number}\" correct? [default: y]: $(tput sgr 0)" tag_correct </dev/tty
        tag_correct=${tag_correct:-$default}
    done

    

fi

openshift_version=`oc version 2>/dev/null|grep "oc v"|cut -d '.' -f2`

operator_files=( 
    "00-namespace.yaml" "01-serviceaccount.yaml"
    "02-alamedaservice.crd.yaml" "03-federatorai-operator.deployment.yaml"
    "04-clusterrole.yaml" "05-clusterrolebinding.yaml"
    "06-role.yaml" "07-rolebinding.yaml"
)

file_folder="/tmp/install-op"

rm -rf $file_folder
mkdir -p $file_folder
current_location=`pwd`
cd $file_folder

for file in "${operator_files[@]}"
do
    echo "Downloading file $file ..."
    if ! curl -sL --fail https://raw.githubusercontent.com/containers-ai/federatorai-operator/${tag_number}/deploy/upstream/${file} -O; then
        echo -e "\n$(tput setaf 1)Abort, download file failed!!!$(tput sgr 0)"
        echo "Please check tag name and network"
        exit 1
    fi
    echo "Done"
done

sed -i "s/ubi:latest/ubi:${tag_number}/g" 03*.yaml

operator_namespace=`cat 00-name*.yaml|grep "name:"|awk '{print $2}'`
echo -e "\n$(tput setaf 2)Starting apply Federator.ai operator yaml files$(tput sgr 0)"
kubectl apply -f .
echo "Processing..."
wait_until_pods_ready 600 20 $operator_namespace 1
echo -e "\n$(tput setaf 6)Install Federator.ai operator $tag_number successfully$(tput sgr 0)"

alamedaservice_example="alamedaservice_sample.yaml"
alamedascaler_example="alamedascaler.yaml"

echo -e "\nDownloading alamedaservice and alamedascaler sample files ..."
if ! curl -sL --fail https://raw.githubusercontent.com/containers-ai/federatorai-operator/${tag_number}/example/${alamedaservice_example} -O; then
    echo -e "\n$(tput setaf 1)Abort, download alamedaservice sample file failed!!!$(tput sgr 0)"
    exit 2
fi

if ! curl -sL --fail https://raw.githubusercontent.com/containers-ai/alameda/${tag_number}/example/samples/nginx/${alamedascaler_example} -O; then
    echo -e "\n$(tput setaf 1)Abort, download alamedascaler sample file failed!!!$(tput sgr 0)"
    exit 3
fi
echo "Done"

sed -i "s/version: latest/version: ${tag_number}/g" ${alamedaservice_example}

echo "========================================"

while [[ "$interactive_install" != "y" ]] && [[ "$interactive_install" != "n" ]]
do
    default="y"
    read -r -p "$(tput setaf 2)Do you want to launch interactive installation of Federator.ai [default: y]: $(tput sgr 0)" interactive_install </dev/tty
    interactive_install=${interactive_install:-$default}
done

if [[ "$interactive_install" == "y" ]]; then

    while [[ "$information_correct" != "y" ]] && [[ "$information_correct" != "Y" ]]
    do
        # init variables
        install_namespace=""
        enable_execution=""
        prometheus_address=""
        storage_type=""
        log_size=""
        data_size=""
        storage_class=""

        default="alameda"
        read -r -p "$(tput setaf 127)Enter the namespace you want to install Federator.ai [default: alameda]: $(tput sgr 0)" install_namespace </dev/tty
        install_namespace=${install_namespace:-$default}
    
        default="y"
        read -r -p "$(tput setaf 127)Do you want to enable execution? [default: y]: $(tput sgr 0): " enable_execution </dev/tty
        enable_execution=${enable_execution:-$default}

        if [[ "$openshift_version" == "11" ]]; then
            default="https://prometheus-k8s.openshift-monitoring:9091"
        elif [[ "$openshift_version" == "9" ]]; then
            default="http://prom-prometheus-operator-prometheus.monitoring.svc:9090"
        else
            default="https://prometheus-k8s.openshift-monitoring:9091"
        fi

        echo "$(tput setaf 127)Enter the prometheus service address"
        read -r -p "[default: ${default}]: $(tput sgr 0)" prometheus_address </dev/tty
        prometheus_address=${prometheus_address:-$default}
        
        while [[ "$storage_type" != "ephemeral" ]] && [[ "$storage_type" != "persistent" ]]
        do
            default="ephemeral"
            echo "$(tput setaf 127)Which storage type you would like to use? ephemeral or persistent?"
            read -r -p "[default: ephemeral]: $(tput sgr 0)" storage_type </dev/tty
            storage_type=${storage_type:-$default}
        done

        if [[ "$storage_type" == "persistent" ]]; then
            default="10"
            read -r -p "$(tput setaf 127)Specify log storage size [ex: 10 for 10GB, default: 10]: $(tput sgr 0)" log_size </dev/tty
            log_size=${log_size:-$default}
            default="10"
            read -r -p "$(tput setaf 127)Specify data storage size [ex: 10 for 10GB, default: 10]: $(tput sgr 0)" data_size </dev/tty
            data_size=${data_size:-$default}
            
            while [[ "$storage_class" == "" ]]
            do
            read -r -p "$(tput setaf 127)Specify storage class name: $(tput sgr 0)" storage_class </dev/tty
            done
            
        fi

        echo -e "\n----------------------------------------"
        echo "install_namespace = $install_namespace"
        if [[ "$enable_execution" == "y" ]]; then
            echo "enable_execution = true"    
        else
            echo "enable_execution = false"
        fi
        echo "prometheus_address = $prometheus_address"
        echo "storage_type = $storage_type"
        if [[ "$storage_type" == "persistent" ]]; then
            echo "log storage size = $log_size GB"
            echo "data storage size = $data_size GB"
            echo "storage class name = $storage_class"
        fi
        echo "----------------------------------------"

        default="y"
        read -r -p "$(tput setaf 2)Is the above information correct [default: y]:$(tput sgr 0)" information_correct </dev/tty
        information_correct=${information_correct:-$default}
    done

    sed -i "s|\bnamespace:.*|namespace: ${install_namespace}|g" ${alamedaservice_example}

    if [[ "$enable_execution" == "y" ]]; then
        sed -i "s/\benableExecution:.*/enableExecution: true/g" ${alamedaservice_example}    
    else
        sed -i "s/\benableExecution:.*/enableExecution: false/g" ${alamedaservice_example}
    fi

    sed -i "s|\bprometheusService:.*|prometheusService: ${prometheus_address}|g" ${alamedaservice_example}
    if [[ "$storage_type" == "persistent" ]]; then
            sed -i '/- usage:/,+10d' ${alamedaservice_example}
            cat >> ${alamedaservice_example} << __EOF__       
    - usage: log 
      type: pvc
      size: ${log_size}Gi
      class: ${storage_class}             
    - usage: data
      type: pvc
      size: ${data_size}Gi
      class: ${storage_class}

__EOF__
    fi
    kubectl create ns $install_namespace &>/dev/null
    kubectl apply -f $alamedaservice_example &>/dev/null
    echo "Processing..."
    wait_until_pods_ready 900 20 $install_namespace 5
    echo -e "$(tput setaf 6)\nInstall Alameda $tag_number successfully$(tput sgr 0)"
    get_grafana_route $install_namespace
    leave_prog
    exit 0
fi

leave_prog
exit 0




