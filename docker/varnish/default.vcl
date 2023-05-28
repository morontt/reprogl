vcl 4.0;

backend default {
    .host = "gopher";
    .port = "8080";
}

sub vcl_backend_response {
    if (beresp.http.Surrogate-Control ~ "ESI/1.0") {
        unset beresp.http.Surrogate-Control;
        set beresp.do_esi = true;
    }
}

sub vcl_deliver {
    unset resp.http.Age;
    unset resp.http.Via;
    unset resp.http.X-Varnish;
}
