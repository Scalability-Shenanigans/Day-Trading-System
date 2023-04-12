import { useEffect, useState } from "react";
import "./App.css";
import TransactionForm from "./components/TransactionForm";
import styled from "styled-components";
import {
  addFunds,
  buyStock,
  commitBuy,
  commitSell,
  getAllTransactionsByUser,
  getBalance,
  getQuote,
  getStocks,
  sellStock,
} from "./requests/requests";
import TransactionsList, {
  TransactionsListItemProps,
} from "./components/TransactionsList";
import LoginAlert from "./components/LoginAlert";
import { canSellStock } from "./utils/utils";

const AddFundsText = styled.h3({
  margin: 0,
});

const AddFundsInput = styled.input({
  marginBottom: 20,
  marginRight: 5,
});

const AppContainer = styled.div`
  display: flex;
  flex-direction: column;
  height: calc(100vh - 120px);
  box-sizing: border-box;
  max-width: 760px;
  margin: 0 auto;
  text-align: center;
  font-family: "Lato", sans-serif;
  background-color: gray;
  padding: 15px;
  border-radius: 10px;
`;

const FormsContainer = styled.div`
  display: flex;
  justify-content: space-between;

  & > * {
    flex: 1;
    margin: 20px 25px;
    box-sizing: border-box;
  }
`;

const TopContainer = styled.div`
  color: white;
`;

const CommitButtonsContainer = styled.div`
  display: flex;
  flex-direction: horizontal;
  justify-content: center;
`;

const CommitButton = styled.button`
  margin: 0px 5px;
`;

function App() {
  const [user, setUser] = useState("");

  const [selectedFunds, setSelectedFunds] = useState<number | null>(null);
  const [funds, setFunds] = useState(0);
  const [transactionCommitted, setTransactionCommitted] = useState(false);
  const [quotePrice, setQuotePrice] = useState(null);

  useEffect(() => {
    if (user !== "") {
      fetchBalance();
      fetchUserTransactions();
    }
  }, [user]);

  const handleLoginSubmit = (username: string, password: string) => {
    console.log("Username:", username);
    console.log("Password:", password);
    setUser(username);
  };

  const fetchUserTransactions = async () => {
    const userTransactions = await getAllTransactionsByUser({
      user: user,
    });

    console.log("the user transactions response is", userTransactions);
    setTransactions(userTransactions as TransactionsListItemProps[]);
  };

  const fetchBalance = async () => {
    const getBalanceResponse = await getBalance({ user: user });
    const balance = parseFloat(getBalanceResponse?.data["balance"].toFixed(2));
    setFunds(balance);
  };

  const fetchQuote = async (user: string, stock: string) => {
    const quote = await getQuote({
      user,
      stock,
    });
    setQuotePrice(quote.toFixed(2));
  };

  const handleQuoteSubmit = async (stock: string) => {
    console.log("Quote submitted", stock);
    stock.length > 0 && (await fetchQuote(user, stock));
  };

  const handleBuySubmit = async (stock: string, amount: number) => {
    console.log("Buy submitted", stock, amount);

    if (stock.length > 0 && funds && funds > 0) {
      if (amount > funds) {
        alert("Not enough funds");
      } else {
        await buyStock({
          user,
          stock,
          amount,
        });
      }
    }

    await fetchUserTransactions();
  };

  const handleSellSubmit = async (stock: string, amount: number) => {
    console.log("Sell submitted", stock, amount);

    if (stock.length > 0 && funds && funds > 0) {
      const userStockHoldings = await getStocks({ user });
      const currentStockPrice = await getQuote({ user, stock });

      console.log("stock holdings are", userStockHoldings);

      if (!canSellStock(userStockHoldings, stock, amount, currentStockPrice)) {
        alert("Insufficient shares to sell");
      } else {
        await sellStock({
          user,
          stock,
          amount,
        });
      }
    }

    await fetchUserTransactions();
  };

  useEffect(() => {
    console.log("transaction committed");

    if (transactionCommitted === true) {
      setTransactionCommitted(false);
    }

    fetchBalance();
    fetchUserTransactions();
  }, [transactionCommitted]);

  const [transactions, setTransactions] = useState<TransactionsListItemProps[]>(
    []
  );

  return user !== "" ? (
    <AppContainer>
      <TopContainer>
        <h1>Daytrading Frontend</h1>
        <h1>Logged in as {user}</h1>

        <h2>Current funds: ${funds}</h2>
        <AddFundsText>Add Funds</AddFundsText>
        <AddFundsInput
          type="number"
          value={selectedFunds ?? ""}
          onChange={(e) =>
            setSelectedFunds(e.target.value ? Number(e.target.value) : 0)
          }
        />
        <button
          onClick={async () => {
            const newBalance = (await addFunds({
              user: user,
              amount: selectedFunds ?? 0,
            })) as unknown as number;
            console.log("newBalance is", newBalance);
            newBalance && setFunds(newBalance);
          }}
        >
          Add funds
        </button>
      </TopContainer>
      <CommitButtonsContainer>
        <CommitButton
          onClick={async () => {
            await commitBuy({
              user: user,
            });
            setTransactionCommitted(true);
          }}
        >
          Commit Buy
        </CommitButton>
        <CommitButton
          onClick={async () => {
            await commitSell({
              user: user,
            });
            setTransactionCommitted(true);
          }}
        >
          Commit Sell
        </CommitButton>
      </CommitButtonsContainer>
      <FormsContainer>
        <TransactionForm
          title="Quote"
          buttonText="Get Quote"
          onSubmit={(stock) => handleQuoteSubmit(stock)}
          showAmount={false}
          quotePrice={quotePrice}
        />
        <TransactionForm
          title="Buy"
          buttonText="Buy"
          onSubmit={(stock, amount) => handleBuySubmit(stock, amount)}
          showAmount={true}
        />
        <TransactionForm
          title="Sell"
          buttonText="Sell"
          onSubmit={(stock, amount) => handleSellSubmit(stock, amount)}
          showAmount={true}
        />
      </FormsContainer>
      <TransactionsList
        transactions={transactions}
        user={user}
        setCommit={setTransactionCommitted}
      />
    </AppContainer>
  ) : (
    <LoginAlert onSubmit={handleLoginSubmit} />
  );
}

export default App;
