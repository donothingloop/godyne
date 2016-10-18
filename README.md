# godyne - Dynamic DNS updater

This project is a dynamic DNS updater with a standard API. It can be used to map your own domains to dynamic IP addresses.
It uses dynamic updates to update the entries of your DNS server.

# Usage

The godyne API server listens on port 33009 by default. This can be changed in the config.json file.

A request may look like this:
```
http://127.0.0.1:33009/api/update?hostname=example.com&myip=8.8.8.8
```

Basic authentication is used. To create a password hash for an user, *htpasswd* can be used like this:
```
htpasswd -n -B admin
```

This will ask you for a password and output an hash like this:
```
test:$2y$05$K/l4v8/iT.kkGGfNPo8/z.8olOqCKknV5rqkyMWvMZnR5X7qKG/UW
```

To add the user to godyne, add the hash to the config.json:
```
"users": {
    "admin": {
        "password": "$2y$05$K/l4v8/iT.kkGGfNPo8/z.8olOqCKknV5rqkyMWvMZnR5X7qKG/UW"
    }
},
```

# Installation

**Please check that you have configured your local GO installation according to this tutorial: https://golang.org/doc/install**

```
go get github.com/donothingloop/godyne
go install github.com/donothingloop/godyne
```

Then create the config directory and copy the sample file:

```
mkdir /etc/godyne/
cp config.sample.js /etc/godyne/config.js
```

Finally copy the systemd service file and enable it:

```
cp godyne.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable godyne.service
systemctl start godyne.service
```

# Configuration

To be useful you will have to configure your DNS server to allow dynamic updates.

At first create the required TSIG keys with *dnssec-keygen*:
On ubuntu enter:

```
sudo apt install bind9utils
dnssec-keygen -a HMAC-SHA512 -b 512 -n HOST example.com.
```

This will generate the files *Kexample.com.+165+22914.key* and *Kexample.com.+165+22914.private*.A

To configure your *bind9* dns server copy the following snippet to your */etc/named.conf* file:

Replace the secret in the quotes with the key from the *Kexample.com.+165+22914.private* file.
Keep the generated files secret!

```
key "example.com." {
        algorithm HMAC-SHA512;
        secret "VfT1btoiuN7QRBGIv26uFyrc5cPOxyhdBq7j8MoPhD/egGsrFJWXFGkN d3JCBusqiO9Z3Rsy2evt625Vjacglw==";
};

```

To enable updates also add this line to your zone definition:

```
allow-update { key "example.com."; };
```

This will allow zone updates from the generated key.
