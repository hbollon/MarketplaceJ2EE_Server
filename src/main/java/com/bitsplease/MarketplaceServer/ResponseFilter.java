package com.bitsplease.MarketplaceServer;

import javax.servlet.*;
import javax.servlet.annotation.WebFilter;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.util.logging.Logger;

@WebFilter(filterName = "ResponseFilter", urlPatterns = {"/*"})
public class ResponseFilter implements Filter {
    private final static Logger log = Logger.getLogger(ResponseFilter.class.getName());

    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain) throws IOException, ServletException {

        HttpServletResponse http = (HttpServletResponse) response;
        http.addHeader("Access-Control-Allow-Origin", "*");
        http.addHeader("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, x-client-key, x-client-token, x-client-secret, Authorization");
        http.addHeader("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT");
        chain.doFilter(request, response);
    }
}