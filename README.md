# KBeauty API

Using headless chromium to automate orders.

# Riman API

products - ## Products - API  https://cart-api.riman.com/api/v2/products?cartType=R&countryCode=US&culture=en-US&isCart=true&repSiteUrl=WeKBeauty


orders - https://cart-api.riman.com/api/v1/orders?mainSiteUrl=2043124962&getEnrollerOrders=&getCustomerOrders=&orderNumber=&shipmentNumber=&trackingNumber=&isRefunded=&paidStatus=true&orderType=&orderLevel=&weChatOrderNumber=&startDate=&endDate=&offset=0&limit=20&orderBy=-mainOrdersPK

shipping - https://cart-api.riman.com/api/v1/orders/{OrderNumber}/shipment-products

cart - https://cart-api.riman.com/api/v1/shopping/72357ede-62c0-4e03-9904-f47be1a7f900
order - post - https://cart-api.riman.com/api/v2/order 
payload [
cartKey: "72357ede-62c0-4e03-9904-f47be1a7f900"
countryCode: "US"
mainId: 51311
mainOrderType: 4
salesCampaignFK: null
]