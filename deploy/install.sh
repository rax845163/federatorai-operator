#!/bin/sh

##################################################################
#
#   This script is created for installing Federator.ai Operator
#
##################################################################

show_usage()
{
    cat << __EOF__

    Usage: $0 [-t Federator.ai_Operator_Tag]

__EOF__
    exit 1
}

while getopts "t:" o; do
    case "${o}" in
        t)
            t_arg=${OPTARG}
            ;;
        *)
            show_usage
            ;;
    esac
done

[ "${t_arg}" = "" ] && show_usage
[ "${t_arg}" != "" ] && tag_number="${t_arg}"

operator_files=( 
    "00-namespace.yaml" "01-serviceaccount.yaml"
    "02-alamedaservice.crd.yaml" "03-federatorai-operator.deployment.yaml"
    "04-clusterrole.yaml" "05-clusterrolebinding.yaml"
    "06-role.yaml" "07-rolebinding.yaml"
)

file_folder="/tmp/install-op"

rm -rf $file_folder
mkdir -p $file_folder
cd $file_folder

for file in "${operator_files[@]}"
do
    echo "Downloading file $file ..."
    if ! curl -sL --fail https://raw.githubusercontent.com/containers-ai/federatorai-operator/${tag_number}/deploy/upstream/${file} -O; then
        echo -e "\nAbort, download file failed!!!"
        echo "Please check tag name and network"
        exit 1
    fi
    echo "Done"
done

sed -i "s/ubi:latest/ubi:${tag_number}/g" 03*.yaml

echo -e "\nStarting apply yaml files"
kubectl apply -f .

echo -e "\nInstall Federator.ai operator $tag_number successfully"

alamedaservice_example="alamedaservice_sample.yaml"
alamedascaler_example="alamedascaler.yaml"

echo -e "\nDownloading alamedaservice and alamedascaler sample files ..."
if ! curl -sL --fail https://raw.githubusercontent.com/containers-ai/federatorai-operator/${tag_number}/example/${alamedaservice_example} -O; then
    echo -e "\nAbort, download alamedaservice sample file failed!!!"
    exit 2
fi

if ! curl -sL --fail https://raw.githubusercontent.com/containers-ai/alameda/${tag_number}/example/samples/nginx/${alamedascaler_example} -O; then
    echo -e "\nAbort, download alamedascaler sample file failed!!!"
    exit 3
fi
echo "Done"

sed -i "s/version: latest/version: ${tag_number}/g" ${alamedaservice_example}

echo -e "\nYAML files are located under $file_folder"
cd - > /dev/null
exit 0

