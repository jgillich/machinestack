address = ":9610"

log_level = "DEBUG"

allow_origins = ["*"]

tls {
  enable = true

  # automatically fetch certificate via letsencrypt
  auto = true

  # or manually set certificate
  cert_file = "/tls.crt"
  key_file = "/tls.key"
}

postgres {
  address   = "localhost:6379"
  username = "postgres"
  password  = ""
  database = "faststack"
}


jwt {
  # secret used by the jwt web token.
  # generate by running openssl rand -base64 32
  secret = ""
}

scheduler {
  name = "local"
}

driver {
  enable = ["lxd"]

  options {
    "lxd.remote" = "unix://"
  }
}
