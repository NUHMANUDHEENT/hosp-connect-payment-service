# Payment Service

The **Payment Service** manages payment processing for appointments and treatments, integrates with Razorpay, and handles notifications using Kafka.

---

## **Features**

### **Payment Processing**
- Integrate with **Razorpay** for secure and seamless transactions.
- Generate invoices and track payment history.

### **Webhook Management**
- Handle Razorpay callbacks to verify payment status and update records.

### **Payment Notifications**
- Trigger real-time payment notifications to users via **Kafka** events, consumed by the Notification Service.

---

## **Technology Stack**
- **Backend:** Go (Golang)
- **Payment Gateway:** Razorpay
- **Event Streaming:** Kafka
- **Database:** PostgreSQL or MongoDB (as per setup)

---

## **How to Run**

### Clone the Repository
```bash
git clone https://github.com/your-username/payment-service.git
cd payment-service
