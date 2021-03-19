package com.bitsplease.MarketplaceServer.model;

import com.bitsplease.MarketplaceServer.RestApplication;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import javax.net.ssl.HttpsURLConnection;
import java.io.*;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.Arrays;

public class Client {
    private int clientId;
    private String firstName;
    private String lastName;
    private String birthday;
    private String nationality;
    private String countryOfResidence;
    private String email;

    public Client() {
    }

    public Client(String firstName, String lastName, String email) {
        this.firstName = firstName;
        this.lastName = lastName;
        this.email = email;
    }

    public int getClientId() {
        return clientId;
    }

    public void setClientId(int clientId) {
        this.clientId = clientId;
    }

    public String getFirstName() {
        return firstName;
    }

    public void setFirstName(String firstName) {
        this.firstName = firstName;
    }

    public String getLastName() {
        return lastName;
    }

    public void setLastName(String lastName) {
        this.lastName = lastName;
    }

    public String getBirthday() {
        return birthday;
    }

    public void setBirthday(String birthday) {
        this.birthday = birthday;
    }

    public String getNationality() {
        return nationality;
    }

    public void setNationality(String nationality) {
        this.nationality = nationality;
    }

    public String getCountryOfResidence() {
        return countryOfResidence;
    }

    public void setCountryOfResidence(String countryOfResidence) {
        this.countryOfResidence = countryOfResidence;
    }

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }

    @Override
    public String toString() {
        return "Client{" +
                "firstName='" + firstName + '\'' +
                ", lastName='" + lastName + '\'' +
                ", birthday='" + birthday + '\'' +
                ", nationality='" + nationality + '\'' +
                ", countryOfResidence='" + countryOfResidence + '\'' +
                ", email='" + email + '\'' +
                '}';
    }

    public static int GetClientId(String email) throws IOException {
        URL url = new URL(RestApplication.MANGOPAY_URL + "/users/");
        System.out.println(url.toString());
        HttpURLConnection con = (HttpURLConnection) url.openConnection();
        con.setRequestMethod("GET");
        con.setRequestProperty("Authorization", RestApplication.AUTHORIZATION_HEADER);
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

    public static int RegisterNewClient(Client client) {
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

            URL registerClientUrl = new URL(RestApplication.MANGOPAY_URL + "/users/natural/");
            HttpsURLConnection con = (HttpsURLConnection)registerClientUrl.openConnection();
            con.setRequestMethod("POST");
            con.setRequestProperty("Authorization", RestApplication.AUTHORIZATION_HEADER);
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

    public static int CreateNewWallet(Client client) {
        try {
            String newWalletRequestJson = new JSONObject()
                    .put("Owners", new JSONArray().put(client.getClientId()))
                    .put("Description", "MarketplaceJ2EE's Wallet")
                    .put("Currency", "EUR")
                    .toString();

            System.out.println("Json: " + newWalletRequestJson);

            URL webCardUrl = new URL(RestApplication.MANGOPAY_URL + "/wallets/");
            HttpsURLConnection con = (HttpsURLConnection)webCardUrl.openConnection();
            con.setRequestMethod("POST");
            con.setRequestProperty("Authorization", RestApplication.AUTHORIZATION_HEADER);
            con.setRequestProperty("Content-Length", String.valueOf(newWalletRequestJson.length()));
            con.setRequestProperty("Content-Type","application/json");
            con.setRequestProperty("User-Agent", "Mozilla/5.0");
            con.setDoOutput(true);
            con.setDoInput(true);

            DataOutputStream output = new DataOutputStream(con.getOutputStream());
            output.writeBytes(newWalletRequestJson);
            output.close();
            DataInputStream input = new DataInputStream( con.getInputStream() );

            String response = "";
            for( int c = input.read(); c != -1; c = input.read() )
                response = response.concat(String.valueOf((char)c));
            input.close();
            System.out.println("Resp Code:" + con.getResponseCode());
            System.out.println("Resp Message:" + con.getResponseMessage());

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
