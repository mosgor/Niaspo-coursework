import React from 'react';

const ProductList = ({ products, addToCart }) => {
  if (!Array.isArray(products)) {
	return <p>No products available</p>;
  }
  return (
    <div className="product-list">
      {products.map((product) => (
        <div key={product.id} className="product">
          <img src={product.image_url} alt={product.name}/>
          <h3>{product.name}</h3>
          <p>{product.description}</p>
          <p>Price: ${product.price}</p>
          <p>Weight: {product.weight}kg</p>
          <button onClick={() => addToCart(product)}>Add to Cart</button>
        </div>
      ))}
    </div>
  );
};

export default ProductList;
