FROM alpine
ADD shopcategory /shopcategory
ENTRYPOINT ["/shopcategory"]