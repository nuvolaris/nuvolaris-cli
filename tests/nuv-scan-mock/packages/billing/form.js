function main() {
    return {
        body: `<html><head>
<link href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.0/css/bootstrap.min.css" 
      rel="stylesheet" id="bootstrap-css">
</head><body>
 <div id="container">
  <div class="row">
   <div class="col-md-8 col-md-offset-2">
    <h4><strong>Pay</strong></h4>
    <form method="POST" action="submit-sendmail">
     <div class="form-group">
       <input type="text" name="credit"
        class="form-control" placeholder="Credit Card Number">
     </div>
     <div class="form-group">
        <input type="text" name="name"
         class="form-control" placeholder="Name">
     </div>
     <div class="form-group">
        <input type="text" name="secret" 
         class="form-control" placeholder="CVC">
     </div>
      <button class="btn btn-default" type="submit" name="button">
        Pay Now
      </button>
    </form>
   </div>
  </div>
 </div>
</body></html>`
    }
}
