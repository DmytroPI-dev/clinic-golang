import Footer from "./components/Footer";
import ImageGallery from "react-image-gallery";
import { useTranslation } from "react-i18next";
import { Container } from "reactstrap";

const images = [];
const importAll = (r) => {
  r.keys().forEach((key) => {
    images.push({
      original: r(key),
      thumbnail: r(key),
      originalAlt: "Villa Dorota rooms",
    });
  });
};

importAll(
  require.context("../static/rooms/", false, /\.(png|jpe?g|svg)$/)
);

const Booking = () => {
  const { t } = useTranslation();
  return (
    <>
      <Container fluid>
        <div className="dormedDescription ms-3">
          <h2 className="text-center mt-3">{t("Willa Dorotka pokoje")}</h2>
          <p>
            {t(
              " W willi posiadamy pokoje jednoosobowe, dwuosobowe oraz pokoje rodzinne."
            )}
          </p>
          <p>{t("W pokojach jednoosobowych znajduje się łoże małżeńskie")}</p>
          <p>
            {t(
              "Pokoje dwuosobowe posiadają dwa pojedyncze łóżka.Pokoje rodzinne mają jedno łoze małżeńskie oraz jedno łóżko pojedyncze"
            )}
          </p>
          <p>
            {t(
              "Wszystkie pokoje z łazienkami (kabina prysznicowa, umywalka, wc, suszarka do włosów, ręczniki) i dostępem do Internetu. Internet bezprzewodowy, również na zewnątrz budynku"
            )}
          </p>
          <p>
            {t(
              "W każdym pokoju znajduje się telewizor, radio, lodówka, czajnik bezprzewodowy i komplet naczyń (talerze, sztućce, szklanki, kieliszki)"
            )}
          </p>
          <p>
            {t(
              "Do dyspozycji gości również: żelazko i deska do prasowania, leżaki, plac zabaw dla dzieci, grill, monitorowany parking"
            )}
          </p>
        </div>

        <ImageGallery
          items={images}
          autoPlay={true}
          thumbnailPosition={"bottom"}
        />
        <div className="dormedDescription ms-3">
          <h3 className="text-center">
            {t("W celu rezerwacji proszę skontaktować się z nami:")}
          </h3>

          <p>{t("kom. 111 11 11 11")}</p>
          <p>{t("kom. 111 11 11 11")}</p>
          <p>{t("kom. 111 11 11 11")}</p>
          <p>{t("tel. 111 11 11 11")}</p>
          <p>
            Email:{" "}
            <a href="mailto: email@email.pl">email@email.pl</a>
          </p>
          <p>
            Email:{" "}
            <a href="mailto: email@gmail.com">email@gmail.com</a>
          </p>

          <h3 className="text-center">
            {t(
              " Nasz system rezerwacyjny wymaga wpłaty zadatku w wysokości do 30% kwoty rezerwacji. Po dokonaniu wpłaty rezerwacja uznawana jest za GWARANTOWANĄ."
            )}
          </h3>
          <h4>{t("Płatności można dokonać:")}</h4>
          <p>{t("wpłacając zadatek na konto firmy:")}</p>
          <p>"Dormed Medical Spa"</p>
          <p>"ul. Newstreet 5"</p>
          <p>"28-100 New city"</p>
          <p>{t("Konto")}</p>
          <p>11 1111 1111 1111 1111 1111</p>
          <p>
            {t(
              "w tytule przelewu proszę podać okres pobytu, ilość osób i obiekt Willa Dorotka lub Dormed Medical SPA"
            )}
          </p>
        </div>
      </Container>
      <Footer />
    </>
  );
};
export default Booking;
