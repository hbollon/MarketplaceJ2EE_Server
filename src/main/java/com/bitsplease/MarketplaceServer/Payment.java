package com.bitsplease.MarketplaceServer;

import com.bitsplease.MarketplaceServer.model.Client;
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
        System.out.println(obj.toString());
        try {
            int userId = getClientId(obj.getClient().getEmail());
            if(userId == -1){
               userId = registerNewClient(obj.getClient());
               if(userId == -1)
                   return "Error during user registration on MangoPay.";
               System.out.println("User successfully registered.");
            }

            String webCardRequestJson = new JSONObject()
                    .put("AuthorId", userId)
                    .put("DebitedFunds", new JSONObject().put("Currency", "EUR").put("Amount", (obj.getProduct().getPrice() + obj.getProduct().getFees()) * 100.0))
                    .put("Fees", new JSONObject().put("Currency", "EUR").put("Amount", 0.0))
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

    public int registerNewClient(Client client) {
        try {
            String clientRegisterRequestJson = new JSONObject()
                    .put("FirstName", client.getFirstName())
                    .put("LastName", client.getLastName())
                    .put("Birthday", 1463496101)
                    .put("Nationality", "FR")
                    .put("CountryOfResidence", "FR")
                    .put("Email", client.getEmail())
                    .toString();

            System.out.println(clientRegisterRequestJson);

            URL registerClientUrl = new URL(MANGOPAY_URL + "/users/natural");
            HttpsURLConnection con = (HttpsURLConnection)registerClientUrl.openConnection();
            con.setRequestMethod("POST");
            con.setRequestProperty("Authorization", AUTHORIZATION_HEADER);
            con.setRequestProperty("Content-Length", String.valueOf(clientRegisterRequestJson.length()));
            con.setRequestProperty("Content-Type","application/json");
            con.setRequestProperty("User-Agent", "Mozilla/5.0");
            con.setDoOutput(true);
            con.setDoInput(true);

            DataOutputStream output = new DataOutputStream(con.getOutputStream());
            output.writeBytes(clientRegisterRequestJson);
            output.close();

            DataInputStream input = new DataInputStream(con.getInputStream());
            String response = "";
            for( int c = input.read(); c != -1; c = input.read() )
                response = response.concat(String.valueOf((char)c));
            input.close();

            if(con.getResponseCode() == HttpURLConnection.HTTP_OK) {
                JSONObject obj = new JSONObject(response);
                return obj.getInt("Id");
            }
            return -1;
        } catch (Exception e) {
            e.printStackTrace();
            return -1;
        }
    }
}