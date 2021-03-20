package com.bitsplease.MarketplaceServer.restservices;

import com.bitsplease.MarketplaceServer.RestApplication;
import com.bitsplease.MarketplaceServer.model.Client;
import com.bitsplease.MarketplaceServer.model.requestObjects.BuyRequest;
import org.json.JSONObject;

import javax.net.ssl.HttpsURLConnection;
import javax.ws.rs.Consumes;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import java.io.DataInputStream;
import java.io.DataOutputStream;
import java.net.URL;

@Path("/mangopay")
public class Payment {
    @Path("/pay")
    @POST
    @Consumes("application/json")
    @Produces("application/json")
    public String pay(
            BuyRequest obj
    ) {
        System.out.println(obj.toString());
        try {
            int userId = Client.GetClientId(obj.getClient().getEmail());
            if(userId == -1){
               userId = Client.RegisterNewClient(obj.getClient());
               if(userId == -1)
                   return "Error during user registration on MangoPay.";
               System.out.println("User successfully registered.");
            }

            String webCardRequestJson = new JSONObject()
                    .put("AuthorId", userId)
                    .put("DebitedFunds", new JSONObject().put("Currency", "EUR").put("Amount", (obj.getProduct().getPrice() + obj.getProduct().getFees()) * 100.0))
                    .put("Fees", new JSONObject().put("Currency", "EUR").put("Amount", 0.0))
                    .put("ReturnURL", "https://hbollon.github.io/MarketplaceJ2EE_Client/")
                    .put("CardType", "CB_VISA_MASTERCARD")
                    .put("CreditedWalletId", Integer.toString(obj.getProduct().getSeller().getWalletId()))
                    .put("Culture", "FR")
                    .toString();

            URL webCardUrl = new URL(RestApplication.MANGOPAY_URL + "/payins/card/web");
            HttpsURLConnection con = (HttpsURLConnection)webCardUrl.openConnection();
            con.setRequestMethod("POST");
            con.setRequestProperty("Authorization", RestApplication.AUTHORIZATION_HEADER);
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
}