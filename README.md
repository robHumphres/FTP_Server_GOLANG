# FTP_Server_GOLANG

This is a golang service that deals as a FTP server. The goal of this is to receive a request a POST request from a user
with the body having file name of "file" and then requires a zip file. 

The API route would be http://localhost:9000/upload

Example of a request:

    public static void transferFile(String fileToTransfer){
        System.out.println("Starting to transfer... " +  fileToTransfer);
        try {
            CloseableHttpClient client = HttpClients.createDefault();
            HttpPost httpPost = new HttpPost("http://localhost:9000/upload");
            MultipartEntityBuilder builder = MultipartEntityBuilder.create();
            builder.addBinaryBody("file", new File(fileToTransfer),
                    ContentType.APPLICATION_JSON, "file.zip");
            HttpEntity multipart = builder.build();

            httpPost.setEntity(multipart);

            CloseableHttpResponse response = client.execute(httpPost);
            assertThat(response.getStatusLine().getStatusCode(), equalTo(200));
            client.close();
        }catch (Exception e){
            System.out.println("There was a fatal error in transfering the file");
        }
    }


The Zipped file will be unzipped and placed in the same location of executable.
