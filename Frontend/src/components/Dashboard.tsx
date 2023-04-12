import React, { useEffect, useState } from "react";
import styled from "styled-components";
import {
  addFunds,
  buyStock,
  commitBuy,
  commitSell,
  getQuote,
  getStocks,
  sellStock,
} from "../requests/requests";
import { canSellStock } from "../utils/utils";
import TransactionForm from "./TransactionForm";

interface DashboardProps {
  user: string;
  funds: number;
  fetchBalance: () => Promise<void>;
  setFunds: React.Dispatch<React.SetStateAction<number>>;
  setTransactionCommitted: React.Dispatch<React.SetStateAction<boolean>>;
  fetchUserTransactions: () => Promise<void>;
}

const AddFundsInput = styled.input`
  margin-bottom: 20px;
  margin-right: 5px;
  background: #2c3e50;
  color: #ecf0f1;
  border: 1px solid #7f8c8d;
  border-radius: 5px;
  padding: 5px;
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

const TopHalfContainer = styled.div`
  background: #34495e;
  padding: 20px;
  border-radius: 10px 10px 0 0;
  margin-bottom: 15px;
`;

const TopContainer = styled.div`
  color: #ecf0f1;
`;

const CommitButtonsContainer = styled.div`
  display: flex;
  flex-direction: horizontal;
  justify-content: center;
  margin-top: 20px;
`;

const CommitButton = styled.button`
  background: #19232e;
  color: #ecf0f1;
  border: none;
  border-radius: 5px;
  padding: 8px 16px;
  margin: 0px 5px;
  cursor: pointer;
  font-weight: bold;

  &:hover {
    background: #2c3e50;
  }
`;

const AddFundsButton = styled.button`
  background: #19232e;
  color: #ecf0f1;
  border: none;
  border-radius: 5px;
  padding: 8px 16px;
  margin: 0px 5px;
  cursor: pointer;
  font-weight: bold;

  &:hover {
    background: #2c3e50;
  }
`;

const Dashboard: React.FC<DashboardProps> = ({
  user,
  funds,
  setFunds,
  fetchBalance,
  fetchUserTransactions,
  setTransactionCommitted,
}) => {
  const [selectedFunds, setSelectedFunds] = useState<number | null>(null);
  const [quotePrice, setQuotePrice] = useState(null);

  const fetchQuote = async (user: string, stock: string) => {
    const quote = await getQuote({
      user,
      stock,
    });
    setQuotePrice(quote.toFixed(2));
  };

  const handleQuoteSubmit = async (stock: string) => {
    stock.length > 0 && (await fetchQuote(user, stock));
  };

  const handleBuySubmit = async (stock: string, amount: number) => {
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
  };

  const handleSellSubmit = async (stock: string, amount: number) => {
    if (stock.length > 0 && funds && funds > 0) {
      const userStockHoldings = await getStocks({ user });
      const currentStockPrice = await getQuote({ user, stock });

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
  };

  useEffect(() => {
    if (user !== "") {
      fetchBalance();
      fetchUserTransactions();
    }
  }, [user]);

  useEffect(() => {
    fetchUserTransactions();
  }, [handleBuySubmit, handleSellSubmit]);

  return (
    <TopHalfContainer>
      <TopContainer>
        <h1>Logged in as {user}</h1>
        <h2>Current funds: ${funds}</h2>
        <AddFundsInput
          type="number"
          value={selectedFunds ?? ""}
          onChange={(e) =>
            setSelectedFunds(e.target.value ? Number(e.target.value) : 0)
          }
        />
        <AddFundsButton
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
        </AddFundsButton>
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
    </TopHalfContainer>
  );
};

export default Dashboard;
