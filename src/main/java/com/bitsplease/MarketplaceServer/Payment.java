package com.bitsplease.MarketplaceServer;

import com.bitsplease.MarketplaceServer.model.Client;
import com.bitsplease.MarketplaceServer.model.Product;
import com.bitsplease.MarketplaceServer.model.requestObjects.BuyRequest;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import javax.net.ssl.*;
import javax.ws.rs.*;
import java.io.*;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.Base64;

@Path("/mangopay")
public class Payment {
    private static final String API_KEY = "fVOz2GjFiOSQdcUocALLaN2AhnOQgQtPv0NDuPQYoDzXKFCpQz";
    private static final String CLIENT_ID = "info802marketplace";
    private static final String AUTHORIZATION_HEADER = "Basic " + Base64.getEncoder().encodeToString((CLIENT_ID + ":" + API_KEY).getBytes());
    private static final String MANGOPAY_URL = "https://api.sandbox.mangopay.com/v2.01/info802marketplace";

    @Path("/pay")
    @POST
    @Consumes("application/json")
    @Produces("application/json")
    public String pay(
            BuyRequest obj
    ) {
        try {
            int userId = getClientId(obj.getClient().getEmail());
            String webCardRequestJson = new JSONObject()
                    .put("AuthorId", userId)
                    .put("DebitedFunds", new JSONObject().put("Currency", "EUR").put("Amount", (obj.getProduct().getPrice() + (obj.getProduct().getWeight() * 0.3)) * 100))
                    .put("Fees", new JSONObject().put("Currency", "EUR").put("Amount", "0"))
                    .put("ReturnURL", "http://localhost:4200")
                    .put("CardType", "CB_VISA_MASTERCARD")
                    .put("CreditedWalletId", "100218510")
                    .put("Culture", "FR")
                    .toString();

            URL webCardUrl = new URL(MANGOPAY_URL + "/payins/card/web");
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

    public int getClientId(String email) throws IOException {
        URL url = new URL(MANGOPAY_URL + "/users/");
        System.out.println(url.toString());
        HttpURLConnection con = (HttpURLConnection) url.openConnection();
        con.setRequestMethod("GET");
        con.setRequestProperty("Authorization", AUTHORIZATION_HEADER);
        con.setRequestProperty("User-Agent", "Mozilla/5.0");

        int responseCode = con.getResponseCode();
        if (responseCode == HttpURLConnection.HTTP_OK) {
            BufferedReader in = new BufferedReader(new InputStreamReader(
                    con.getInputStream()));

            String inputLine;
            StringBuilder response = new StringBuilder();
            while ((inputLine = in.readLine()) != null) {
                response.append(inputLine);
            }
            in.close();
            System.out.println(response.toString());
            try {
                JSONArray responseJson = new JSONArray(response.toString());
                for (int i = 0; i < responseJson.length(); i++) {
                    JSONObject obj = responseJson.getJSONObject(i);
                    if(obj.getString("Email").equals(email)) {
                        return obj.getInt("Id");
                    }
                }
            } catch (JSONException e) {
                e.printStackTrace();
            }
            return -1;
        } else {
            System.out.println("GET request not worked");
        }
        return -1;
    }

}