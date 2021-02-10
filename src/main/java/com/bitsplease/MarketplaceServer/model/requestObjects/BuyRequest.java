package com.bitsplease.MarketplaceServer.model.requestObjects;

import com.bitsplease.MarketplaceServer.model.Client;
import com.bitsplease.MarketplaceServer.model.Product;

public class BuyRequest {
    private Client client;
    private Product product;

    public BuyRequest() {
    }

    public Client getClient() {
        return client;
    }

    public void setClient(Client client) {
        this.client = client;
    }

    public Product getProduct() {
        return product;
    }

    public void setProduct(Product product) {
        this.product = product;
    }

    @Override
    public String toString() {
        return "BuyRequest{" +
                "client=" + client.toString() +
                ", product=" + product.toString() +
                '}';
    }
}
