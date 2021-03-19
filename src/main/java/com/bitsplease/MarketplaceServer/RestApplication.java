package com.bitsplease.MarketplaceServer;

import javax.net.ssl.*;
import javax.ws.rs.ApplicationPath;
import javax.ws.rs.core.Application;
import java.security.cert.X509Certificate;
import java.util.Base64;

@ApplicationPath("/rest")
public class RestApplication extends Application {
    public static final String API_KEY = "fVOz2GjFiOSQdcUocALLaN2AhnOQgQtPv0NDuPQYoDzXKFCpQz";
    public static final String CLIENT_ID = "info802marketplace";
    public static final String AUTHORIZATION_HEADER = "Basic " + Base64.getEncoder().encodeToString((CLIENT_ID + ":" + API_KEY).getBytes());
    public static final String MANGOPAY_URL = "https://api.sandbox.mangopay.com/v2.01/info802marketplace";

    public RestApplication() {
        super();
        try {
            this.FixHttpsCert();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    /**
     * Used to get rid of: "javax.net.ssl.SSLHandshakeException: PKIX path building failed:
     * sun.security.provider.certpath.SunCertPathBuilderException:
     * unable to find valid certification path to requested target" with mangopay sandbox api
     *
     * It override the SSL checking process
     * @throws Exception
     */
    public void FixHttpsCert() throws Exception {
        TrustManager[] trustAllCerts = new TrustManager[] {
                new X509TrustManager() {
                    public java.security.cert.X509Certificate[] getAcceptedIssuers() {
                        return null;
                    }
                    public void checkClientTrusted(X509Certificate[] certs, String authType) {  }
                    public void checkServerTrusted(X509Certificate[] certs, String authType) {  }
                }
        };

        SSLContext sc = SSLContext.getInstance("SSL");
        sc.init(null, trustAllCerts, new java.security.SecureRandom());
        HttpsURLConnection.setDefaultSSLSocketFactory(sc.getSocketFactory());

        // Create all-trusting host name verifier
        HostnameVerifier allHostsValid = new HostnameVerifier() {
            public boolean verify(String hostname, SSLSession session) {
                return true;
            }
        };
        // Install the all-trusting host verifier
        HttpsURLConnection.setDefaultHostnameVerifier(allHostsValid);
    }

}