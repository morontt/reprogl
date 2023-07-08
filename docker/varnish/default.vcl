vcl 4.0;

import std;

backend default {
    .host = "gopher";
    .port = "8080";
}

sub vcl_recv {
    if (req.method == "POST" && req.url == "/purge-cache") {
        if (req.http.x-ban-token == std.getenv("BAN_TOKEN")) {
            std.ban("req.http.host ~ .*");
        } else {
            return (synth(403, "Forbidden"));
        }
    }

    if (req.http.Cookie) {
        if (req.http.Cookie ~ ".*(^|;| )?session=([a-zA-Z0-9\-_=]+)( |;|$)?.*") {
            set req.http.X-Varnish-Session = regsub(req.http.Cookie, ".*(^|;| )?session=([a-zA-Z0-9\-_=]+)( |;|$)?.*", "\2");
        }

        if (!(req.url ~ "^/login")) {
            unset req.http.Cookie;
        }
    }
}

sub vcl_backend_response {
    if (beresp.http.Surrogate-Control ~ "ESI/1.0") {
        unset beresp.http.Surrogate-Control;
        set beresp.do_esi = true;
    }
}

sub vcl_hash {
    hash_data(req.url);
    if (req.http.X-Varnish-Session && (
        req.url == "/_fragment/auth-navigation" ||
        req.url ~ "^/_fragment/comments/")
    ) {
        hash_data(req.http.X-Varnish-Session);
    }

    return (lookup);
}

sub vcl_deliver {
    unset resp.http.Age;
    unset resp.http.Via;
    unset resp.http.X-Varnish;
}
