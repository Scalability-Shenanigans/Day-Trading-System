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

interface userOnlyProps {
  user: string;
}

interface getQuoteProps {
  user: string;
  stock: string;
}

interface TransactionData {
  type: string;
  date: string;
  asset: string;
  amount: number;
  timestamp: number;
  isCommitted: boolean;
}

type Transaction = {
  Transaction_ID: number;
  Stock: string;
  Is_Buy: boolean;
  Amount: number;
  Price: number;
  User: string;
  Timestamp: number;
};

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
    return response;
  } catch (error) {
    console.log("the error", error);
    return null;
  }
}

async function commitBuy({ user }: userOnlyProps) {
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

async function commitSell({ user }: userOnlyProps) {
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

async function getBalance({ user }: userOnlyProps) {
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

async function getAllTransactionsByUser({
  user,
}: userOnlyProps): Promise<TransactionData[]> {
  const data = {
    user,
  };

  try {
    const response = await axios.post(
      `${transaction_server_url}/allTransactionsByUser`,
      JSON.stringify(data),
      {
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    const pendingTransactions = response.data["pending_transactions"];
    const finishedTransactions = response.data["finished_transactions"];

    console.log("pending transactions are", pendingTransactions);
    console.log("finished transactions are", finishedTransactions);

    const formattedPendingTransactions =
      pendingTransactions?.map((transaction: Transaction) => ({
        type: transaction.Is_Buy ? "Buy" : "Sell",
        date: new Date(transaction.Timestamp).toISOString().slice(0, 10),
        asset: transaction.Stock,
        amount: transaction.Amount,
        user: transaction.User,
        timestamp: transaction.Timestamp,
        isCommitted: false,
      })) ?? [];

    // Map finished transactions to the desired format
    const formattedFinishedTransactions =
      finishedTransactions?.map((transaction: Transaction) => ({
        type: transaction.Is_Buy ? "Buy" : "Sell",
        date: new Date(transaction.Timestamp).toISOString().slice(0, 10),
        asset: transaction.Stock,
        amount: transaction.Amount,
        user: transaction.User,
        timestamp: transaction.Timestamp,
        isCommitted: true,
      })) ?? [];

    const allTransactions = formattedPendingTransactions
      ?.concat(formattedFinishedTransactions)
      ?.sort(
        (a: TransactionData, b: TransactionData) => a.timestamp - b.timestamp
      )
      .map(({ timestamp, ...rest }: TransactionData) => rest);

    return allTransactions;
  } catch (error) {
    console.log("the error", error);
    return [];
  }
}

async function getQuote({ user, stock }: getQuoteProps) {
  const data = {
    user,
    stock,
  };

  try {
    const response = await axios.post(
      `${transaction_server_url}/quote`,
      JSON.stringify(data),
      {
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    console.log("the response is", response);

    return response.data["price"];
  } catch (error) {
    console.log("the error", error);
    return null;
  }
}

export {
  addFunds,
  getQuote,
  buyStock,
  commitBuy,
  sellStock,
  commitSell,
  getBalance,
  getAllTransactionsByUser,
};
