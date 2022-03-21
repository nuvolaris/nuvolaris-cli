def main(args):
    name = args.get("name", "stranger")
    return { 
      "body": f"hello {name}!\n"
    } 
