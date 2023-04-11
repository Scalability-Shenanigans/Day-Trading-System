import axios from "axios";

const transaction_server_url = "http://localhost:5100";

interface addFundsProps {
  user: string;
  amount: number;
}

interface buyAndSellStockProps {
  user: string;
  stock: string;
  amount: number;
}

interface commitBuySellGetBalanceProps {
  user: string;
}

async function addFunds({ user, amount }: addFundsProps) {
  const data = {
    user,
    amount,
  };

  try {
    const response = await axios.post(
      `${transaction_server_url}/add`,
      JSON.stringify(data),
      {
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    console.log("the response is", response);
    return response.data["balance"];
  } catch (error) {
    console.log("the error", error);
    return 0;
  }
}

async function buyStock({ user, stock, amount }: buyAndSellStockProps) {
  const data = {
    user,
    stock,
    amount,
  };

  try {
    const response = await axios.post(
      `${transaction_server_url}/buy`,
      JSON.stringify(data),
      {
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    console.log("the response is", response);
    return response.status;
  } catch (error) {
    console.log("the error", error);
    return 0;
  }
}

async function commitBuy({ user }: commitBuySellGetBalanceProps) {
  const data = {
    user,
  };

  try {
    const response = await axios.post(
      `${transaction_server_url}/commitBuy`,
      JSON.stringify(data),
      {
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    console.log("the response is", response);

    return response;
  } catch (error) {
    console.log("the error", error);
    return null;
  }
}

async function sellStock({ user, stock, amount }: buyAndSellStockProps) {
  const data = {
    user,
    stock,
    amount,
  };

  try {
    const response = await axios.post(
      `${transaction_server_url}/sell`,
      JSON.stringify(data),
      {
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    console.log("the response is", response);
    return response;
  } catch (error) {
    console.log("the error", error);
    return null;
  }
}

async function commitSell({ user }: commitBuySellGetBalanceProps) {
  const data = {
    user,
  };

  try {
    const response = await axios.post(
      `${transaction_server_url}/commitSell`,
      JSON.stringify(data),
      {
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    console.log("the response is", response);

    return response;
  } catch (error) {
    console.log("the error", error);
    return null;
  }
}

async function getBalance({ user }: commitBuySellGetBalanceProps) {
  const data = {
    user,
  };

  try {
    const response = await axios.post(
      `${transaction_server_url}/getBalance`,
      JSON.stringify(data),
      {
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    console.log("the response is", response);

    return response;
  } catch (error) {
    console.log("the error", error);
    return null;
  }
}

export { addFunds, buyStock, commitBuy, sellStock, commitSell, getBalance };
