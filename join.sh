#!/bin/bash
set -e

CACHE=0
if [ $# -ne 3 ] ; then
	if [ $# -eq 4 ] && [ $4 = 'CACHE' ] ; then
		CACHE=1
	else
		echo "USAGE: $0 hostname domain username ['CACHE']" >&2
		echo "Example: $0 iwg-gen0 d3catur.net rbmason" >&2
		exit 1
	fi
fi

HOST=$1
DOMAIN=$2
USER=$3

set -x

hostname $HOST
echo "127.0.0.1    $HOST.$DOMAIN      $HOST      localhost" > /etc/hosts
echo "$HOST" > /etc/hostname

apt-get install realmd adcli sssd libsss-sudo crudini samba-common-bin sssd-tools -y
mkdir -p /var/lib/samba/private
systemctl enable sssd

# Get Domain Controller
DCURL=`host -t SRV _ldap._tcp.dc._msdcs.$DOMAIN | cut '-d ' -f8 | sed 's/.$//'`

crudini --set /etc/systemd/timesyncd.conf Time Servers $DCURL
timedatectl set-ntp true
# Let time sync
sleep 5

realm discover $DOMAIN
realm join --user=$USER $DOMAIN
crudini --set /etc/sssd/sssd.conf sssd default_domain_suffix $DOMAIN
systemctl start sssd

# If not found in common-session, add pam_mkhomedir.so to pam common-session so
# domain users will have home directories.
grep pam_mkhomedirs.so /etc/pam.d/common-session > /dev/null || \
    echo 'session required pam_mkhomedir.so' >> /etc/pam.d/common-session

SUDO_FNAME=/etc/sudoers.d/admins_`echo $DOMAIN | cut -d. -f1`
echo "%domain\ admins@$DOMAIN ALL=(root) ALL" >> $SUDO_FNAME
chmod 440 $SUDO_FNAME

[ $CACHE -eq 1 ] && apt-get install nss-updatedb libnss-db libpam-ccreds nscd -y

set +x

echo SUCCESS - you must restart in order to complete the join.

