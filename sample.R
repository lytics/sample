
sampler <- function(x, n = 1e5, replace = FALSE) {
  p <- length(x)
  m <- matrix(NA, nrow = n, ncol = p)
  ip <- 1:p
  for(i in 1:n) {
    m[i,] <- sample(ip, replace = replace, prob = x)
  }
  return(m)
}

sampler.summary <- function(m) {
  for(i in 1:ncol(m)) {
    print(table(m[,i]) / nrow(m))
  }
}
