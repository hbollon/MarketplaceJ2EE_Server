package com.bitsplease.MarketplaceServer;

import com.bitsplease.MarketplaceServer.model.Product;
import org.json.JSONObject;

import javax.net.ssl.*;
import javax.ws.rs.*;
import java.io.DataInputStream;
import java.io.DataOutputStream;
import java.net.URL;
import java.security.cert.X509Certificate;
import java.util.Base64;

@Path("/mangopay")
public class Payment {
    private static final String API_KEY = "fVOz2GjFiOSQdcUocALLaN2AhnOQgQtPv0NDuPQYoDzXKFCpQz";
    private static final String CLIENT_ID = "info802marketplace";
    private static final String AUTHORIZATION_HEADER = "Basic " + Base64.getEncoder().encodeToString((CLIENT_ID + ":" + API_KEY).getBytes());

    @Path("/pay")
    @POST
    @Consumes("application/json")
    @Produces("application/json")
    public String pay(
            Product product
    ) {
        try {
            final String httpsURL = "https://api.sandbox.mangopay.com/v2.01/info802marketplace";
            this.FixHttpsCert();

            int userId = (int)(Math.random()*10000)+1;
            String webCardRequestJson = new JSONObject()
                    .put("AuthorId", 100218276)
                    .put("DebitedFunds", new JSONObject().put("Currency", "EUR").put("Amount", (product.getPrice() + (product.getWeight() * 0.3)) * 100))
                    .put("Fees", new JSONObject().put("Currency", "EUR").put("Amount", "0"))
                    .put("ReturnURL", "http://localhost:4200")
                    .put("CardType", "CB_VISA_MASTERCARD")
                    .put("CreditedWalletId", "100218510")
                    .put("Culture", "FR")
                    .toString();

            URL webCardUrl = new URL(httpsURL+"/payins/card/web");
            HttpsURLConnection con = (HttpsURLConnection)webCardUrl.openConnection();
            con.setRequestMethod("POST");
            con.setRequestProperty("Authorization", AUTHORIZATION_HEADER);
            con.setRequestProperty("Content-Length", String.valueOf(webCardRequestJson.length()));
            con.setRequestProperty("Content-Type","application/json");
            con.setRequestProperty("User-Agent", "Mozilla/5.0");
            con.setDoOutput(true);
            con.setDoInput(true);

            DataOutputStream output = new DataOutputStream(con.getOutputStream());
            output.writeBytes(webCardRequestJson);
            output.close();
            DataInputStream input = new DataInputStream( con.getInputStream() );

            String response = "";
            for( int c = input.read(); c != -1; c = input.read() )
                response = response.concat(String.valueOf((char)c));
            input.close();

            System.out.println("Resp Code:" + con.getResponseCode());
            System.out.println("Resp Message:" + con.getResponseMessage());

            return response;
        } catch (Exception e) {
            e.printStackTrace();
            return e.toString();
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