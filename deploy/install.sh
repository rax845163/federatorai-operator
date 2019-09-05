#!/bin/sh

#################################################################################################################
#
#   This script is created for installing Federator.ai Operator
#
#   1. Interactive Mode
#      Usage: ./install.sh
#
#   2. Silent Mode - Persistent storage
#      Usage: ./install.sh -t v4.2.260 -n federatorai -e y -p https://prometheus-k8s.openshift-monitoring:9091 \
#                   -s persistent -l 11 -d 12 -c managed-nfs-storage
#
#   3. Silent Mode - Ephemeral storage
#      Usage: ./install.sh -t v4.2.260 -n federatorai -e y -p https://prometheus-k8s.openshift-monitoring:9091 \
#                   -s ephemeral
#
#   -t = tag_number
#   -n = install_namespace
#   -e = enable_execution
#   -p = prometheus_address
#   -s = storage_type
#   -l = log_size
#   -d = data_size
#   -c = storage_class
#################################################################################################################

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

  wait_pod_creating=1
  for ((i=0; i<$period; i+=$interval)); do

    if [[ "$wait_pod_creating" = "1" ]]; then
        # check if pods created
        if [[ "`kubectl get po -n $namespace 2>/dev/null|wc -l`" -ge "$target_pod_number" ]]; then
            wait_pod_creating=0
        else
            echo "Waiting for pods in namespace $namespace to be created..."
        fi
    else
        # check if pods running
        if pods_ready $namespace; then
            echo -e "\nAll $namespace pods are ready."
            return 0
        fi
        echo "Waiting for pods in namespace $namespace to be ready..."
    fi

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

while getopts "t:n:e:p:s:l:d:c:" o; do
    case "${o}" in
        t)
            t_arg=${OPTARG}
            ;;
        n)
            n_arg=${OPTARG}
            ;;
        e)
            e_arg=${OPTARG}
            ;;
        p)
            p_arg=${OPTARG}
            ;;
        s)
            s_arg=${OPTARG}
            ;;
        l)
            l_arg=${OPTARG}
            ;;
        d)
            d_arg=${OPTARG}
            ;;
        c)
            c_arg=${OPTARG}
            ;;
        *)
            echo "Warning! wrong paramter, ignore it."
            ;;
    esac
done

[ "${t_arg}" = "" ] && silent_mode_disabled="y"
[ "${n_arg}" = "" ] && silent_mode_disabled="y"
[ "${e_arg}" = "" ] && silent_mode_disabled="y"
[ "${p_arg}" = "" ] && silent_mode_disabled="y"
[ "${s_arg}" = "" ] && silent_mode_disabled="y"
[ "${s_arg}" = "persistent" ] && [ "${l_arg}" = "" ] && silent_mode_disabled="y"
[ "${s_arg}" = "persistent" ] && [ "${d_arg}" = "" ] && silent_mode_disabled="y"
[ "${s_arg}" = "persistent" ] && [ "${c_arg}" = "" ] && silent_mode_disabled="y"

[ "${t_arg}" != "" ] && tag_number="${t_arg}"
[ "${n_arg}" != "" ] && install_namespace="${n_arg}"
[ "${e_arg}" != "" ] && enable_execution="${e_arg}"
[ "${p_arg}" != "" ] && prometheus_address="${p_arg}"
[ "${s_arg}" != "" ] && storage_type="${s_arg}"
[ "${l_arg}" != "" ] && log_size="${l_arg}"
[ "${d_arg}" != "" ] && data_size="${d_arg}"
[ "${c_arg}" != "" ] && storage_class="${c_arg}"

if [ "$silent_mode_disabled" = "y" ];then

    while [[ "$info_correct" != "y" ]] && [[ "$info_correct" != "Y" ]]
    do
        # init variables
        install_namespace=""
        tag_number=""

        read -r -p "$(tput setaf 2)Please input Federator.ai Operator tag:$(tput sgr 0) " tag_number </dev/tty

        default="federatorai"
        read -r -p "$(tput setaf 2)Enter the namespace you want to install Federator.ai [default: federatorai]: $(tput sgr 0)" install_namespace </dev/tty
        install_namespace=${install_namespace:-$default}

        echo -e "\n----------------------------------------"
        echo "tag_number = $tag_number"
        echo "install_namespace = $install_namespace"
        echo "----------------------------------------"

        default="y"
        read -r -p "$(tput setaf 2)Is the above information correct? [default: y]: $(tput sgr 0)" info_correct </dev/tty
        info_correct=${info_correct:-$default}
    done
else
    echo -e "\n----------------------------------------"
    echo "tag_number=$tag_number"
    echo "install_namespace=$install_namespace"
    echo "enable_execution=$enable_execution"
    echo "prometheus_address=$prometheus_address"
    echo "storage_type=$storage_type"
    echo "log_size=$log_size"
    echo "data_size=$data_size"
    echo "storage_class=$storage_class"
    echo -e "----------------------------------------\n"
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

# Modify federator.ai operator yaml(s)
# for tag
sed -i "s/ubi:latest/ubi:${tag_number}/g" 03*.yaml
# for namespace
sed -i "s/name: federatorai/name: ${install_namespace}/g" 00*.yaml
sed -i "s/namespace: federatorai/namespace: ${install_namespace}/g" 01*.yaml 03*.yaml 05*.yaml 06*.yaml 07*.yaml

echo -e "\n$(tput setaf 2)Starting apply Federator.ai operator yaml files$(tput sgr 0)"
for yaml_fn in `ls [0-9]*.yaml | sort -n`; do
    echo "Applying ${yaml_fn}..."
    kubectl apply -f ${yaml_fn}
    if [ "$?" != "0" ]; then
        echo -e "\n$(tput setaf 1)Error in applying yaml file ${yaml_fn}.$(tput sgr 0)"
        exit 5
    fi
done
wait_until_pods_ready 600 30 $install_namespace 1
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

if [ "$silent_mode_disabled" = "y" ];then

    while [[ "$install_alameda" != "y" ]] && [[ "$install_alameda" != "n" ]]
    do
        default="y"
        read -r -p "$(tput setaf 2)Do you want to launch interactive installation of Federator.ai [default: y]: $(tput sgr 0)" install_alameda </dev/tty
        install_alameda=${install_alameda:-$default}
    done

    if [[ "$install_alameda" == "y" ]]; then

        while [[ "$information_correct" != "y" ]] && [[ "$information_correct" != "Y" ]]
        do
            # init variables
            enable_execution=""
            prometheus_address=""
            storage_type=""
            log_size=""
            data_size=""
            storage_class=""

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
    fi
else
    install_alameda="y"
fi

if [[ "$install_alameda" == "y" ]]; then
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

    kubectl apply -f $alamedaservice_example &>/dev/null
    echo "Processing..."
    wait_until_pods_ready 900 60 $install_namespace 5
    echo -e "$(tput setaf 6)\nInstall Alameda $tag_number successfully$(tput sgr 0)"
    get_grafana_route $install_namespace
    leave_prog
    exit 0
fi

leave_prog
exit 0
