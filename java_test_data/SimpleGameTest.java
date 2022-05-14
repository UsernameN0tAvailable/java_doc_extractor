package snakes;

import org.junit.jupiter.api.Test;
import snakes.squares.AlarmSquare;
import snakes.squares.InstantLoseSquare;
import snakes.squares.SleepSquare;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

public class SimpleGameTest {

    private Player jack;
    private Player jill;

    private Game game;

    public Game newGame() {
        jack = new Player("Jack");
        jill = new Player("Jill");
        Player[] args = { jack, jill };
        game = new Game(15, args);
        game.setSquareToLadder(2, 4);
        game.setSquareToLadder(6, 2);
        game.setSquareToSnake(11, -6);
        game.setSquare(3, new AlarmSquare(game, 3));
        game.setSquare(4, new SleepSquare(game, 4));
        game.setSquare(7, new AlarmSquare(game, 7));
        game.setSquare(10, new InstantLoseSquare(game, 10));
        assertTrue(game.notOver());
        assertTrue(game.firstSquare().isOccupied());
        assertEquals(1, jack.position());
        assertEquals(1, jill.position());
        assertEquals(jack, game.currentPlayer());
        return game;
    }

    public Game setAnotherGameState() {
        jack = new Player("Jack");
        jill = new Player("Jill");
        Player[] args = { jack, jill };
        game = new Game(15, args);
        game.setSquareToLadder(2, 4);
        game.setSquareToLadder(6, 2);
        game.setSquareToSnake(11, -6);
        game.setSquare(3, new AlarmSquare(game, 3));
        game.setSquare(4, new SleepSquare(game, 4));
        game.setSquare(10, new InstantLoseSquare(game, 10));
        assertTrue(game.notOver());
        assertTrue(game.firstSquare().isOccupied());
        assertEquals(1, jack.position());
        assertEquals(1, jill.position());
        assertEquals(jack, game.currentPlayer());
        return game;
    }

    @Test
    public void initialStrings() {
        newGame();
        assertEquals("Jack", jack.toString());
        assertEquals("Jill", jill.toString());
        assertEquals("[1<Jack><Jill>]", game.firstSquare().toString());
        assertEquals("[2->6]", game.getSquare(2).toString());
        assertEquals("[5<-11]", game.getSquare(11).toString());
        assertEquals("[1<Jack><Jill>][2->6][3 (ALARM)][4 (Sleep)][5][6->8][7 (ALARM)][8][9][10 (LOSE)][5<-11][12][13][14][15]", game.toString());
    }

    @Test
    public void playGame() {
        newGame();
        game.movePlayer(4); // Jack moves
        assertTrue(game.notOver());
        assertEquals(5, jack.position());
        assertEquals(1, jill.position());
        assertEquals(jill, game.currentPlayer());
        assertEquals("[1<Jill>]", game.firstSquare().toString());
        assertEquals("[5<Jack>]", game.getSquare(5).toString());
        game.movePlayer(5); // Jill moves, lands on ladder 6 -> 8
        assertTrue(game.notOver());
        assertEquals(5, jack.position());
        assertEquals(8, jill.position());
        assertEquals(jack, game.currentPlayer());
        game.movePlayer(4); // Jack's turn
        assertEquals(9, jack.position());
        assertEquals(8, jill.position());
        game.movePlayer(6); // Jill moves
        game.movePlayer(2); // Jack moves
        game.movePlayer(1); // Jill wins
        assertTrue(game.isOver());
        assertEquals(jill, game.winner());
    }

    @Test
    public void playerOutOfBoard() {
        setAnotherGameState();
        assertEquals(jack, game.currentPlayer());

        game.movePlayer(6); // Jack moves
        assertTrue(game.notOver());
        assertEquals(7, jack.position());
        assertEquals(1, jill.position());
        assertEquals(jill, game.currentPlayer());
        assertEquals("[1<Jill>]", game.firstSquare().toString());
        assertEquals("[7<Jack>]", game.getSquare(7).toString());

        game.movePlayer(6); // Jill can not move as the square is already occupied
        assertTrue(game.notOver());
        assertEquals(7, jack.position());
        assertEquals(1, jill.position());
        assertEquals(jack, game.currentPlayer()); // It is Jack turn now

        game.movePlayer(6); // Jack move from square 7 to square 13
        assertEquals(13, jack.position());
        assertEquals(1, jill.position());

        game.movePlayer(6); // Jill moves
        assertEquals(7, jill.position());
        assertEquals(13, jack.position());

        game.movePlayer(3); // Jack moves 3, should not be out of board, either Jack should be on the 14 square or move to the first square
        assertEquals(14, jack.position());
        assertTrue(game.notOver());
    }
}
