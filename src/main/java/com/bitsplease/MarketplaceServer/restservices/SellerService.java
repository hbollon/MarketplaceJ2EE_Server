package com.bitsplease.MarketplaceServer.restservices;

import com.bitsplease.MarketplaceServer.model.Client;
import com.bitsplease.MarketplaceServer.model.Seller;

import javax.ws.rs.Consumes;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;

@Path("/seller")
public class SellerService {
    @Path("/register")
    @POST
    @Consumes("application/json")
    @Produces("application/json")
    public String register(
            Seller obj
    ) {
        Client sellerAccount = new Client(obj.getFirstName(), obj.getLastName(), obj.getEmail());
        try {
            int userId = Client.GetClientId(sellerAccount.getEmail());
            if(userId == -1){
                System.out.println("Client registration");
                userId = Client.RegisterNewClient(sellerAccount);
                if(userId == -1)
                    return "Error during user registration on MangoPay.";
                System.out.println("User successfully registered.");
            }
            sellerAccount.setClientId(userId);
            System.out.println("Client id: " + userId);

            return Integer.toString(Client.CreateNewWallet(sellerAccount));
        } catch (Exception e) {
            e.printStackTrace();
            return e.toString();
        }
    }
}
